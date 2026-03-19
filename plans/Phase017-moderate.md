# Phase 015: 콘텐츠 모더레이션 프레임워크 — pkg/moderate

## 목표

toulmin 기반 콘텐츠 모더레이션 프레임워크를 `pkg/moderate`에 구현한다.
"이 콘텐츠를 게시할 것인가"를 defeats graph로 판정한다.
차단 규칙 + 예외(뉴스 인용, 교육 목적, 신뢰 사용자)를 defeat edge로 선언적 처리한다.
EvaluateTrace가 모더레이션 로그이자 이의 제기 대응 근거가 된다.

## 배경

### 현재 문제

1. **모더레이션 로직이 if-else로 얽힌다**: 혐오 표현 차단 + 뉴스 인용 예외 + 교육 목적 예외 + 풍자 예외... 차단 규칙마다 예외가 있고, 예외의 예외가 있다
2. **규칙 추가/변경이 위험하다**: 새 차단 규칙을 넣으면 기존 예외와 충돌하는지 전체 코드를 읽어야 한다. 사이드이펙트가 빈번하다
3. **판정 근거 추적이 별도 작업이다**: "왜 차단됐는가"를 사용자에게 설명하려면 모더레이션 로그를 따로 구축해야 한다. 이의 제기 대응에 필수인데 구현 부담이 크다
4. **외부 서비스 비용**: AWS Rekognition, Perspective API 등 판정 API 호출 비용이 콘텐츠 양에 비례하여 증가한다

### toulmin이 해결하는 것

- 차단 규칙 = rebuttal (Go 함수, 1-2 depth)
- 예외 = defeat edge (뉴스 인용, 교육 목적 등)
- qualifier = 심각도 (1.0 = 무조건 차단, 0.3 = 경미한 위반)
- 판정 근거 = EvaluateTrace (모더레이션 로그 내장)
- AI 모델 연동 = rule 함수 안에서 호출 (프레임워크는 결과만 받음)

### claim/ground 분리 원칙

toulmin의 `(claim any, ground any)` 시그니처가 프레임워크 확장성의 핵심이다.

- **claim = 뭘 판정하나**: 모더레이션 프레임워크에서 claim은 콘텐츠 자체 (본문, 미디어)
- **ground = 판정 재료**: ground는 판정에 필요한 컨텍스트 (작성자, 채널, 맥락)

프레임워크는 Moderator 구조와 판정 흐름을 제공하고, **도메인 데이터는 ground로 사용자가 주입한다.**

| 역할 | 모더레이션 프레임워크에서 |
|---|---|
| claim | Content (Body, MediaURLs, ContentType) |
| ground | ContentContext (Author, Channel, Metadata) |
| rule 함수 | claim/ground에서 조건 하나만 판단 (1-2 depth) |
| graph | rule 간 관계 선언 (defeat = 차단 예외) |
| qualifier | 위반 심각도 (1.0 = 무조건, 0.3 = 경미) |
| verdict | 게시 허용/차단 판정 |

이전 프레임워크들과의 차이: **rule 함수 안에서 외부 AI 모델을 호출할 수 있다.** ContainsHateSpeech가 LLM/분류 모델을 호출하고 결과를 bool로 반환하면, 프레임워크는 그 결과를 graph에서 다른 rule과 함께 판정한다. AI 판정과 규칙 기반 판정이 같은 graph에 공존한다.

## 핵심 설계

### Content

```go
// pkg/moderate/content.go
type Content struct {
    Body        string
    MediaURLs   []string
    ContentType string   // "text", "image", "video"
    Metadata    map[string]any
}
```

### ContentContext

```go
// pkg/moderate/content_context.go
type ContentContext struct {
    Author   *Author
    Channel  *Channel
    Metadata map[string]any
}

// pkg/moderate/author.go
type Author struct {
    ID        string
    Verified  bool
    PostCount int
    TrustScore float64  // 0.0 ~ 1.0
}

// pkg/moderate/channel.go
type Channel struct {
    ID       string
    Type     string   // "general", "news", "adult", "education"
    AgeGated bool
}
```

### 범용 rule 함수

```go
// pkg/moderate/rule_is_verified_user.go
func IsVerifiedUser(claim any, ground any) (bool, any) {
    ctx := ground.(*ContentContext)
    return ctx.Author.Verified, nil
}

// pkg/moderate/rule_is_trusted_user.go
func IsTrustedUser(minScore float64) toulmin.RuleFunc {
    return func(claim any, ground any) (bool, any) {
        ctx := ground.(*ContentContext)
        return ctx.Author.TrustScore >= minScore, nil
    }
}

// pkg/moderate/rule_contains_hate_speech.go
// Classifier 인터페이스를 받아 AI 모델 연동 가능
func ContainsHateSpeech(classifier Classifier) toulmin.RuleFunc {
    return func(claim any, ground any) (bool, any) {
        content := claim.(*Content)
        score := classifier.Predict(content.Body, "hate_speech")
        return score > 0.8, score  // evidence로 점수 반환
    }
}

// pkg/moderate/rule_contains_spam.go
func ContainsSpam(classifier Classifier) toulmin.RuleFunc {
    return func(claim any, ground any) (bool, any) {
        content := claim.(*Content)
        score := classifier.Predict(content.Body, "spam")
        return score > 0.7, score
    }
}

// pkg/moderate/rule_contains_nsfw.go
func ContainsNSFW(classifier Classifier) toulmin.RuleFunc {
    return func(claim any, ground any) (bool, any) {
        content := claim.(*Content)
        score := classifier.Predict(content.Body, "nsfw")
        return score > 0.8, score
    }
}

// pkg/moderate/rule_is_news_context.go
func IsNewsContext(claim any, ground any) (bool, any) {
    ctx := ground.(*ContentContext)
    return ctx.Channel.Type == "news", nil
}

// pkg/moderate/rule_is_educational.go
func IsEducational(claim any, ground any) (bool, any) {
    ctx := ground.(*ContentContext)
    return ctx.Channel.Type == "education", nil
}

// pkg/moderate/rule_is_adult_channel.go
func IsAdultChannel(claim any, ground any) (bool, any) {
    ctx := ground.(*ContentContext)
    return ctx.Channel.AgeGated, nil
}

// pkg/moderate/rule_has_min_posts.go
func HasMinPosts(min int) toulmin.RuleFunc {
    return func(claim any, ground any) (bool, any) {
        ctx := ground.(*ContentContext)
        return ctx.Author.PostCount >= min, nil
    }
}
```

### Classifier 인터페이스

```go
// pkg/moderate/classifier.go
type Classifier interface {
    // Predict — 콘텐츠를 분류하여 0.0~1.0 점수 반환
    // category: "hate_speech", "spam", "nsfw" 등
    Predict(text string, category string) float64
}
```

사용자가 구현체를 주입한다. LLM, 분류 모델, 키워드 매칭, 외부 API 등 어떤 방식이든 Classifier 인터페이스만 만족하면 된다.

### Moderator — 모더레이션 판정 실행

```go
// pkg/moderate/moderator.go
type Moderator struct {
    graph *toulmin.Graph
}

func NewModerator(g *toulmin.Graph) *Moderator

// Review — 콘텐츠 판정
func (m *Moderator) Review(content *Content, ctx *ContentContext) (*ReviewResult, error)
```

### ReviewResult

```go
// pkg/moderate/review_result.go
type ReviewResult struct {
    Allowed bool
    Verdict float64
    Action  Action              // Allow, Flag, Block
    Trace   []toulmin.TraceEntry
}

// pkg/moderate/action.go
type Action string

const (
    ActionAllow Action = "allow"   // verdict > 0.3
    ActionFlag  Action = "flag"    // 0.0 < verdict <= 0.3 (수동 검토 필요)
    ActionBlock Action = "block"   // verdict <= 0.0 (undecided 포함)
)
```

verdict를 3단계로 매핑한다. 단순 허용/차단이 아니라 **수동 검토(flag)** 영역이 존재한다. qualifier를 심각도로 사용하므로 verdict가 0 근처면 "판단이 애매한 콘텐츠"로 플래그한다.

### 사용 예시

**주의**: 클로저 rule은 변수에 저장 후 재사용해야 한다. Rebuttal만으로는 공격이 일어나지 않으며 반드시 Defeat edge를 선언해야 한다. 예외를 처리하는 rule은 Defeater로 등록해야 한다.

```go
classifier := myClassifier{}  // Classifier 인터페이스 구현체

// 클로저는 변수에 저장 후 재사용
hateSpeech := moderate.ContainsHateSpeech(classifier)
spam := moderate.ContainsSpam(classifier)
nsfw := moderate.ContainsNSFW(classifier)
trustedUser := moderate.IsTrustedUser(0.9)

g := toulmin.NewGraph("post:publish").
    Warrant(moderate.IsVerifiedUser, 1.0).
    Rebuttal(hateSpeech, 1.0).
    Rebuttal(spam, 0.8).
    Rebuttal(nsfw, 1.0).
    Defeater(moderate.IsNewsContext, 1.0).          // 예외 rule은 Defeater로 등록
    Defeater(moderate.IsEducational, 1.0).
    Defeater(trustedUser, 1.0).
    Defeater(moderate.IsAdultChannel, 1.0).
    Defeat(hateSpeech, moderate.IsVerifiedUser).    // Rebuttal → Warrant 공격 edge 필수
    Defeat(spam, moderate.IsVerifiedUser).
    Defeat(nsfw, moderate.IsVerifiedUser).
    Defeat(moderate.IsNewsContext, hateSpeech).     // Defeater → Rebuttal 예외 처리
    Defeat(moderate.IsEducational, hateSpeech).
    Defeat(trustedUser, spam).
    Defeat(moderate.IsAdultChannel, nsfw)

mod := moderate.NewModerator(g)

content := &moderate.Content{
    Body:        "기사 인용: ...",
    ContentType: "text",
}

ctx := &moderate.ContentContext{
    Author:  &moderate.Author{ID: "user-1", Verified: true, TrustScore: 0.95},
    Channel: &moderate.Channel{ID: "ch-news", Type: "news", AgeGated: false},
}

result, err := mod.Review(content, ctx)
// result.Action: "allow" (뉴스 채널이므로 혐오 표현 차단이 defeat됨)
// result.Trace: ContainsHateSpeech=true(score=0.92), IsNewsContext=true → defeat
```

### Gin 미들웨어

```go
// pkg/moderate/gin_moderate.go
// Guard — 콘텐츠 제출 엔드포인트에 모더레이션 판정 적용
func Guard(m *Moderator, extractor ContentExtractor) gin.HandlerFunc

// ContentExtractor — gin.Context에서 Content와 ContentContext를 추출
type ContentExtractor func(*gin.Context) (*Content, *ContentContext)
```

```go
r.POST("/posts", moderate.Guard(mod, extractPost), createPostHandler)
// Block → 403, Flag → 202 (수동 검토 대기), Allow → c.Next()
```

## 범위

### 포함

1. **Content, ContentContext 구조체**: 모더레이션 판정에 필요한 콘텐츠/컨텍스트
2. **Author, Channel 구조체**: 작성자/채널 모델
3. **Classifier 인터페이스**: AI 모델 연동 추상화
4. **범용 rule 함수**: IsVerifiedUser, IsTrustedUser, ContainsHateSpeech, ContainsSpam, ContainsNSFW, IsNewsContext, IsEducational, IsAdultChannel, HasMinPosts
5. **Moderator**: 모더레이션 판정 실행 (Review)
6. **ReviewResult, Action**: 판정 결과 (allow/flag/block) + trace
7. **Gin 미들웨어**: Guard
8. **테스트**: rule 함수 단위 테스트, Moderator 통합 테스트

### 제외

- Classifier 구현체 (인터페이스만 정의, LLM/모델은 사용자 책임)
- 이미지/비디오 분석 (텍스트 분류만, 미디어는 Classifier 구현체에서 처리)
- 모더레이션 이력 퍼시스턴스 (DB 저장은 사용자 책임)
- 이의 제기 워크플로우 — pkg/approve와 조합하여 사용 가능
- 대시보드/관리 UI

## 산출물

```
pkg/
  moderate/
    content.go                     — Content 구조체
    content_context.go             — ContentContext 구조체
    author.go                      — Author 구조체
    channel.go                     — Channel 구조체
    classifier.go                  — Classifier 인터페이스
    action.go                      — Action 상수 (allow/flag/block)
    rule_is_verified_user.go       — IsVerifiedUser
    rule_is_trusted_user.go        — IsTrustedUser (클로저)
    rule_contains_hate_speech.go   — ContainsHateSpeech (클로저, Classifier)
    rule_contains_spam.go          — ContainsSpam (클로저, Classifier)
    rule_contains_nsfw.go          — ContainsNSFW (클로저, Classifier)
    rule_is_news_context.go        — IsNewsContext
    rule_is_educational.go         — IsEducational
    rule_is_adult_channel.go       — IsAdultChannel
    rule_has_min_posts.go          — HasMinPosts (클로저)
    moderator.go                   — Moderator (NewModerator, Review)
    review_result.go               — ReviewResult 구조체
    gin_moderate.go                — Guard (Gin 미들웨어)
    rule_test.go                   — rule 함수 단위 테스트
    moderator_test.go              — Moderator 통합 테스트
    gin_moderate_test.go           — Gin 미들웨어 테스트
```

## 단계

### Step 1: 구조체 및 인터페이스 정의

- Content, ContentContext, Author, Channel
- Classifier 인터페이스
- ReviewResult, Action

### Step 2: rule 함수 구현

- 각 rule 함수를 파일 하나에 하나씩 구현 (filefunc 규칙 준수)
- 각 함수는 1-2 depth 유지
- 클로저 rule: IsTrustedUser, ContainsHateSpeech, ContainsSpam, ContainsNSFW, HasMinPosts
- ContainsHateSpeech 등은 Classifier를 받아 AI 판정 결과를 bool로 변환

### Step 3: Moderator 구현

- NewModerator: graph를 받아 Moderator 생성
- Review: graph.EvaluateTrace(content, ctx) → verdict를 Action으로 매핑 → ReviewResult
- Action 매핑: verdict >= 0.3 → allow, 0.0 <= verdict < 0.3 → flag, verdict < 0.0 → block

### Step 4: Gin 미들웨어 구현

- Guard: ContentExtractor로 요청에서 Content/ContentContext 추출 → Review → Action에 따라 응답
- Block → 403, Flag → 202 (수동 검토 대기), Allow → c.Next()

### Step 5: 테스트

- rule 함수 단위 테스트: mock Classifier로 각 조건별 true/false
- Moderator 통합 테스트:
  - 혐오 표현 없음 → allow
  - 혐오 표현 감지 → block
  - 혐오 표현 + 뉴스 채널 defeat → allow
  - 스팸 + 신뢰 사용자 defeat → allow
  - NSFW + 성인 채널 defeat → allow
  - 경미한 위반 (낮은 qualifier) → flag
  - trace에 각 rule의 activated/defeated + evidence(점수) 포함
- Gin 미들웨어 테스트: block → 403, flag → 202, allow → next

### Step 6: 전체 테스트 PASS 확인

- `go test ./...` 전체 PASS 확인

## 검증 기준

1. ContainsHateSpeech 등 rule 함수가 Classifier 결과를 올바르게 변환한다
2. evidence에 Classifier 점수가 포함된다
3. Moderator.Review가 verdict를 Action(allow/flag/block)으로 정확히 매핑한다
4. defeat edge가 맥락 예외를 정확히 처리한다 (뉴스 인용 → 혐오 표현 차단 무시)
5. qualifier 심각도에 따라 flag 영역이 동작한다
6. ReviewResult.Trace가 모더레이션 로그로 사용 가능하다 (각 rule 판정 + 점수)
7. Gin Guard가 Action에 따라 올바른 HTTP 상태 코드를 반환한다
8. 전체 테스트 PASS

## 의존성

- Phase 001-009: toulmin 코어 (NewGraph, Evaluate, EvaluateTrace)
