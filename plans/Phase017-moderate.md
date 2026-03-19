# Phase 017: 콘텐츠 모더레이션 프레임워크 — pkg/moderate (구현 완료)

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
- 판정 근거 = EvaluateTrace (모더레이션 로그 내장, backing 값까지 추적)
- AI 모델 연동 = rule 함수 안에서 호출 (프레임워크는 결과만 받음)

### claim/ground/backing 분리 원칙

toulmin의 `func(claim any, ground any, backing any) (bool, any)` 시그니처가 프레임워크 확장성의 핵심이다.

- **claim = 뭘 판정하나**: 모더레이션 프레임워크에서 claim은 콘텐츠 자체 (본문, 미디어)
- **ground = 판정 재료**: ground는 판정에 필요한 컨텍스트 (작성자, 채널, 맥락). 평가 시점에 요청마다 달라진다
- **backing = 규칙의 판정 기준**: backing은 graph 구성 시 고정되는 값 (Classifier, 임계값, 최소 게시 수 등). 엔진이 관리하는 명시적 값이다

프레임워크는 Moderator 구조와 판정 흐름을 제공하고, **도메인 데이터는 ground로, 판정 기준은 backing으로 사용자가 주입한다.**

| 역할 | 모더레이션 프레임워크에서 | 예시 |
|---|---|---|
| claim | Content (Body, MediaURLs, ContentType) | 게시글 본문, 이미지 URL |
| ground | ContentContext (Author, Channel, Metadata) | 작성자 정보, 채널 종류 |
| backing | 규칙의 판정 기준 (graph 구성 시 고정) | Classifier, 임계값 0.8, 최소 게시 수 10 |
| rule 함수 | claim/ground/backing에서 조건 하나만 판단 (1-2 depth) | 순수 함수, 클로저 없음 |
| graph | rule 간 관계 선언 (defeat = 차단 예외) | |
| qualifier | 위반 심각도 (1.0 = 무조건, 0.3 = 경미) | |
| verdict | 게시 허용/차단 판정 | <= 0 block, 0 < v <= 0.3 flag, > 0.3 allow |

이전 프레임워크들과의 차이: **rule 함수 안에서 외부 AI 모델을 호출할 수 있다.** ContainsHateSpeech가 backing으로 받은 Classifier를 사용해 LLM/분류 모델을 호출하고 결과를 bool로 반환하면, 프레임워크는 그 결과를 graph에서 다른 rule과 함께 판정한다. AI 판정과 규칙 기반 판정이 같은 graph에 공존한다.

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

모든 rule 함수는 `func(claim any, ground any, backing any) (bool, any)` 시그니처의 순수 함수이다. 클로저를 사용하지 않는다. 판정 기준은 backing으로 전달받는다.

```go
// pkg/moderate/rule_is_verified_user.go
// backing: nil (판정 기준 없음 — ground의 Verified 필드만 사용)
func IsVerifiedUser(claim any, ground any, backing any) (bool, any) {
    ctx := ground.(*ContentContext)
    return ctx.Author.Verified, nil
}

// pkg/moderate/rule_is_trusted_user.go
// backing: float64 (최소 신뢰 점수)
func IsTrustedUser(claim any, ground any, backing any) (bool, any) {
    ctx := ground.(*ContentContext)
    minScore := backing.(float64)
    return ctx.Author.TrustScore >= minScore, nil
}

// pkg/moderate/rule_contains_hate_speech.go
// backing: Classifier (AI 분류 모델)
func ContainsHateSpeech(claim any, ground any, backing any) (bool, any) {
    content := claim.(*Content)
    classifier := backing.(Classifier)
    score := classifier.Predict(content.Body, "hate_speech")
    return score > 0.8, score  // evidence로 점수 반환
}

// pkg/moderate/rule_contains_spam.go
// backing: Classifier (AI 분류 모델)
func ContainsSpam(claim any, ground any, backing any) (bool, any) {
    content := claim.(*Content)
    classifier := backing.(Classifier)
    score := classifier.Predict(content.Body, "spam")
    return score > 0.7, score
}

// pkg/moderate/rule_contains_nsfw.go
// backing: Classifier (AI 분류 모델)
func ContainsNSFW(claim any, ground any, backing any) (bool, any) {
    content := claim.(*Content)
    classifier := backing.(Classifier)
    score := classifier.Predict(content.Body, "nsfw")
    return score > 0.8, score
}

// pkg/moderate/rule_is_news_context.go
// backing: nil
func IsNewsContext(claim any, ground any, backing any) (bool, any) {
    ctx := ground.(*ContentContext)
    return ctx.Channel.Type == "news", nil
}

// pkg/moderate/rule_is_educational.go
// backing: nil
func IsEducational(claim any, ground any, backing any) (bool, any) {
    ctx := ground.(*ContentContext)
    return ctx.Channel.Type == "education", nil
}

// pkg/moderate/rule_is_adult_channel.go
// backing: nil
func IsAdultChannel(claim any, ground any, backing any) (bool, any) {
    ctx := ground.(*ContentContext)
    return ctx.Channel.AgeGated, nil
}

// pkg/moderate/rule_has_min_posts.go
// backing: int (최소 게시 수)
func HasMinPosts(claim any, ground any, backing any) (bool, any) {
    ctx := ground.(*ContentContext)
    min := backing.(int)
    return ctx.Author.PostCount >= min, nil
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

사용자가 구현체를 주입한다. LLM, 분류 모델, 키워드 매칭, 외부 API 등 어떤 방식이든 Classifier 인터페이스만 만족하면 된다. Classifier는 backing으로 rule 함수에 전달된다.

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
    ActionFlag  Action = "flag"    // 0 < verdict <= 0.3 (수동 검토 필요)
    ActionBlock Action = "block"   // verdict <= 0 (undecided 포함)
)
```

verdict를 3단계로 매핑한다. 단순 허용/차단이 아니라 **수동 검토(flag)** 영역이 존재한다. qualifier를 심각도로 사용하므로 verdict가 0 근처면 "판단이 애매한 콘텐츠"로 플래그한다.

### 사용 예시

체이닝 없이 각 rule을 `*Rule` 참조로 받고, `Defeat`는 참조로 관계만 선언한다.

```go
classifier := myClassifier{}  // Classifier 인터페이스 구현체

g := toulmin.NewGraph("post:publish")

// 정의 — 각 rule이 *Rule 참조를 반환
verified := g.Warrant(moderate.IsVerifiedUser, nil, 1.0)
hate     := g.Rebuttal(moderate.ContainsHateSpeech, classifier, 1.0)
spam     := g.Rebuttal(moderate.ContainsSpam, classifier, 0.8)
nsfw     := g.Rebuttal(moderate.ContainsNSFW, classifier, 1.0)
news     := g.Defeater(moderate.IsNewsContext, nil, 1.0)
edu      := g.Defeater(moderate.IsEducational, nil, 1.0)
trusted  := g.Defeater(moderate.IsTrustedUser, 0.9, 1.0)
adult    := g.Defeater(moderate.IsAdultChannel, nil, 1.0)

// 관계 — *Rule 참조로 가리킴
g.Defeat(hate, verified)
g.Defeat(spam, verified)
g.Defeat(nsfw, verified)
g.Defeat(news, hate)
g.Defeat(edu, hate)
g.Defeat(trusted, spam)
g.Defeat(adult, nsfw)

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
// result.Trace: ContainsHateSpeech=true(score=0.92, backing=classifier), IsNewsContext=true → defeat
```

Defeat는 `*Rule` 참조 두 개만 받는다. backing이 있는 rule이든 없는 rule이든 동일하게 `g.Defeat(from, to)`로 선언한다. DefeatWith는 없다.

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
4. **범용 rule 함수**: IsVerifiedUser, IsTrustedUser, ContainsHateSpeech, ContainsSpam, ContainsNSFW, IsNewsContext, IsEducational, IsAdultChannel, HasMinPosts — 모두 순수 함수, backing으로 판정 기준 전달
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
    rule_is_verified_user.go       — IsVerifiedUser (backing: nil)
    rule_is_trusted_user.go        — IsTrustedUser (backing: float64, 최소 신뢰 점수)
    rule_contains_hate_speech.go   — ContainsHateSpeech (backing: Classifier)
    rule_contains_spam.go          — ContainsSpam (backing: Classifier)
    rule_contains_nsfw.go          — ContainsNSFW (backing: Classifier)
    rule_is_news_context.go        — IsNewsContext (backing: nil)
    rule_is_educational.go         — IsEducational (backing: nil)
    rule_is_adult_channel.go       — IsAdultChannel (backing: nil)
    rule_has_min_posts.go          — HasMinPosts (backing: int, 최소 게시 수)
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
- 모든 rule 함수는 `func(claim any, ground any, backing any) (bool, any)` 시그니처의 순수 함수
- backing이 필요한 rule: IsTrustedUser(float64), ContainsHateSpeech(Classifier), ContainsSpam(Classifier), ContainsNSFW(Classifier), HasMinPosts(int)
- backing이 nil인 rule: IsVerifiedUser, IsNewsContext, IsEducational, IsAdultChannel

### Step 3: Moderator 구현

- NewModerator: graph를 받아 Moderator 생성
- Review: graph.EvaluateTrace(content, ctx) → verdict를 Action으로 매핑 → ReviewResult
- Action 매핑: verdict > 0.3 → allow, 0 < verdict <= 0.3 → flag, verdict <= 0 → block

### Step 4: Gin 미들웨어 구현

- Guard: ContentExtractor로 요청에서 Content/ContentContext 추출 → Review → Action에 따라 응답
- Block → 403, Flag → 202 (수동 검토 대기), Allow → c.Next()

### Step 5: 테스트

- rule 함수 단위 테스트: mock Classifier를 backing으로 전달하여 각 조건별 true/false
- Moderator 통합 테스트:
  - 혐오 표현 없음 → allow
  - 혐오 표현 감지 → block
  - 혐오 표현 + 뉴스 채널 defeat → allow
  - 스팸 + 신뢰 사용자 defeat → allow
  - NSFW + 성인 채널 defeat → allow
  - 경미한 위반 (낮은 qualifier) → flag
  - trace에 각 rule의 activated/defeated + evidence(점수) + backing 포함
- Gin 미들웨어 테스트: block → 403, flag → 202, allow → next

### Step 6: 전체 테스트 PASS 확인

- `go test ./...` 전체 PASS 확인

## 검증 기준

1. 모든 rule 함수가 `func(claim any, ground any, backing any) (bool, any)` 순수 함수이다 (클로저 없음)
2. ContainsHateSpeech 등 rule 함수가 backing으로 받은 Classifier 결과를 올바르게 변환한다
3. evidence에 Classifier 점수가 포함된다
4. EvaluateTrace의 TraceEntry에 backing 값이 포함된다 (Classifier, 임계값 등)
5. Moderator.Review가 verdict를 Action(allow/flag/block)으로 정확히 매핑한다 (verdict <= 0 block, 0 < verdict <= 0.3 flag, > 0.3 allow)
6. defeat edge가 맥락 예외를 정확히 처리한다 (뉴스 인용 → 혐오 표현 차단 무시)
7. qualifier 심각도에 따라 flag 영역이 동작한다
8. ReviewResult.Trace가 모더레이션 로그로 사용 가능하다 (각 rule 판정 + 점수 + backing)
9. Defeat가 `*Rule` 참조 두 개만 받아 정상 동작한다
10. Gin Guard가 Action에 따라 올바른 HTTP 상태 코드를 반환한다
11. 전체 테스트 PASS

## 의존성

- Phase 001-009: toulmin 코어 (NewGraph, Evaluate, EvaluateTrace)
- Phase 010: backing 일급 시민 (3-element 시그니처, Warrant/Rebuttal/Defeater backing 인자, TraceEntry backing)
- Phase 012: Rule 참조 반환 + 체이닝 제거 (Warrant/Rebuttal/Defeater → *Rule, Defeat(*Rule, *Rule), DefeatWith 제거, GraphBuilder → Graph)
