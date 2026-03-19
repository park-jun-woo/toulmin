# Phase 014: 피처 플래그 프레임워크 — pkg/feature

## 목표

toulmin 기반 피처 플래그 프레임워크를 `pkg/feature`에 구현한다.
"이 기능을 이 사용자에게 켤 것인가"를 defeats graph로 판정한다.
LaunchDarkly 같은 외부 SaaS 없이 Go 바이너리 안에서 피처 토글, 점진적 롤아웃, 예외 처리를 선언적으로 수행한다.

## 배경

### 현재 문제

1. **피처 플래그가 if-else로 관리된다**: `if isBeta || (isKR && rand < 0.3) && !isLegacy` 같은 분기가 코드 곳곳에 박힌다
2. **외부 SaaS 의존**: LaunchDarkly(월 $175~), Split.io, ConfigCat 등 피처 플래그만으로 월 수십만원. SDK 의존성, 네트워크 호출, 장애 전파 위험
3. **예외 처리가 세그먼트 설정으로 흩어진다**: 대시보드 UI에서 세그먼트를 만들고 규칙을 설정하는데, 규칙 간 충돌이나 예외의 예외는 표현이 어렵다
4. **판정 근거 추적이 유료**: "왜 이 사용자에게 이 기능이 꺼졌는가"를 추적하려면 LaunchDarkly 상위 플랜이 필요하다

### toulmin이 해결하는 것

- 피처 활성화 조건 = rule 함수 (1-2 depth)
- 조건 충돌/예외 = defeats graph (명시적 선언)
- 점진적 롤아웃 = qualifier (0.3 = 30% 롤아웃)
- 판정 근거 = EvaluateTrace (내장, 무료)
- 외부 의존성 = 없음 (순수 Go)

### claim/ground 분리 원칙

toulmin의 `(claim any, ground any)` 시그니처가 프레임워크 확장성의 핵심이다.

- **claim = 뭘 판정하나**: 피처 프레임워크에서 claim은 피처 이름
- **ground = 판정 재료**: ground는 피처 판정에 필요한 컨텍스트 (사용자 속성, 환경)

프레임워크는 Flag 구조와 판정 흐름을 제공하고, **도메인 데이터는 ground로 사용자가 주입한다.**

| 역할 | 피처 프레임워크에서 |
|---|---|
| claim | 피처 이름 (string) |
| ground | UserContext (ID, Role, Region, Attributes) |
| rule 함수 | ground에서 조건 하나만 판단 (1-2 depth) |
| graph | rule 간 관계 선언 (defeat = 예외) |
| qualifier | 롤아웃 비율 (0.3 = 30%) 또는 1.0 (전체) |
| verdict | 활성화/비활성화 판정 |

이전 프레임워크들과의 차이: price(Phase 013)에서 qualifier가 할인율이었다면, 여기서는 **qualifier가 롤아웃 비율**이 된다. 같은 메커니즘이 도메인에 따라 다른 의미를 갖는다.

## 핵심 설계

### UserContext

```go
// pkg/feature/user_context.go
type UserContext struct {
    ID         string
    Role       string
    Region     string
    Attributes map[string]any
}
```

### 범용 rule 함수

```go
// pkg/feature/rule_is_beta_user.go
func IsBetaUser(claim any, ground any) (bool, any) {
    ctx := ground.(*UserContext)
    beta, _ := ctx.Attributes["beta"].(bool)
    return beta, nil
}

// pkg/feature/rule_is_internal_staff.go
func IsInternalStaff(claim any, ground any) (bool, any) {
    ctx := ground.(*UserContext)
    return ctx.Role == "internal", nil
}

// pkg/feature/rule_is_region.go
func IsRegion(region string) toulmin.RuleFunc {
    return func(claim any, ground any) (bool, any) {
        ctx := ground.(*UserContext)
        return ctx.Region == region, nil
    }
}

// pkg/feature/rule_has_attribute.go
func HasAttribute(key string, value any) toulmin.RuleFunc {
    return func(claim any, ground any) (bool, any) {
        ctx := ground.(*UserContext)
        return ctx.Attributes[key] == value, nil
    }
}

// pkg/feature/rule_is_legacy_browser.go
func IsLegacyBrowser(claim any, ground any) (bool, any) {
    ctx := ground.(*UserContext)
    legacy, _ := ctx.Attributes["legacy_browser"].(bool)
    return legacy, nil
}

// pkg/feature/rule_is_user_in_percentage.go
// 사용자 ID 기반 결정론적 해시로 롤아웃 비율 판정 (rand 아님)
func IsUserInPercentage(pct float64) toulmin.RuleFunc {
    return func(claim any, ground any) (bool, any) {
        ctx := ground.(*UserContext)
        return hashPercentage(ctx.ID) < pct, nil
    }
}
```

### Flags — 피처 플래그 레지스트리

```go
// pkg/feature/flags.go
type Flags struct {
    features map[string]*toulmin.Graph
}

func NewFlags() *Flags

// Register — 피처 graph 등록
func (f *Flags) Register(name string, g *toulmin.Graph) *Flags

// IsEnabled — 피처 활성화 여부 판정
func (f *Flags) IsEnabled(name string, ctx *UserContext) (bool, error)

// Evaluate — 피처 판정 + verdict 값
func (f *Flags) Evaluate(name string, ctx *UserContext) (float64, error)

// EvaluateTrace — 피처 판정 + 근거
func (f *Flags) EvaluateTrace(name string, ctx *UserContext) (*FeatureResult, error)

// List — 사용자에 대해 활성화된 전체 피처 목록
func (f *Flags) List(ctx *UserContext) ([]string, error)
```

### FeatureResult

```go
// pkg/feature/feature_result.go
type FeatureResult struct {
    Name    string
    Enabled bool
    Verdict float64
    Trace   []toulmin.TraceEntry
}
```

### hashPercentage

```go
// pkg/feature/hash_percentage.go
// hashPercentage — 사용자 ID를 결정론적 해시로 0.0~1.0 비율에 매핑
// rand가 아니므로 같은 사용자는 항상 같은 결과
func hashPercentage(userID string) float64
```

### 사용 예시

**주의**: 클로저 rule은 변수에 저장 후 재사용해야 한다. Rebuttal만으로는 공격이 일어나지 않으며 반드시 Defeat edge를 선언해야 한다. 예외를 처리하는 rule은 Defeater로 등록해야 한다.

```go
flags := feature.NewFlags()

// 클로저는 변수에 저장 후 재사용
regionKR := feature.IsRegion("KR")
pct10 := feature.IsUserInPercentage(0.1)

flags.Register("dark-mode",
    toulmin.NewGraph("feature:dark-mode").
        Warrant(feature.IsBetaUser, 1.0).
        Warrant(regionKR, 0.3).                          // 한국 30% 롤아웃
        Rebuttal(feature.IsLegacyBrowser, 1.0).
        Defeater(feature.IsInternalStaff, 1.0).          // 예외 rule은 Defeater로 등록
        Defeat(feature.IsLegacyBrowser, feature.IsBetaUser). // Rebuttal → Warrant 공격 edge 필수
        Defeat(feature.IsInternalStaff, feature.IsLegacyBrowser), // Defeater → Rebuttal 예외 처리
)

flags.Register("new-checkout",
    toulmin.NewGraph("feature:new-checkout").
        Warrant(pct10, 1.0).                             // 전체 10% 롤아웃
        Warrant(feature.IsInternalStaff, 1.0),           // 내부 직원은 전원
)

ctx := &feature.UserContext{
    ID:     "user-123",
    Role:   "internal",
    Region: "KR",
    Attributes: map[string]any{
        "beta":           true,
        "legacy_browser": false,
    },
}

// 단순 on/off
enabled, _ := flags.IsEnabled("dark-mode", ctx)

// 판정 근거 추적
result, _ := flags.EvaluateTrace("dark-mode", ctx)
// result.Enabled: true
// result.Trace: IsBetaUser=true, IsRegion(KR)=true, IsLegacyBrowser=false

// 전체 활성 피처 목록
active, _ := flags.List(ctx)
// ["dark-mode", "new-checkout"]
```

### Gin 미들웨어

```go
// pkg/feature/gin_feature.go
// Require — 특정 피처가 활성화된 사용자만 접근 허용
func Require(f *Flags, name string, ctxBuilder FeatureContextFunc) gin.HandlerFunc

// Inject — 활성 피처 목록을 gin.Context에 주입 (핸들러에서 조회 가능)
func Inject(f *Flags, ctxBuilder FeatureContextFunc) gin.HandlerFunc
```

```go
r := gin.Default()
r.GET("/checkout/v2", feature.Require(flags, "new-checkout", buildCtx), handler)
r.Use(feature.Inject(flags, buildCtx))  // 모든 핸들러에서 활성 피처 조회 가능
```

## 범위

### 포함

1. **UserContext 구조체**: 피처 판정에 필요한 사용자 컨텍스트
2. **범용 rule 함수**: IsBetaUser, IsInternalStaff, IsRegion, HasAttribute, IsLegacyBrowser, IsUserInPercentage
3. **hashPercentage**: 결정론적 사용자 ID 해시 (롤아웃용)
4. **Flags**: 피처 레지스트리 (Register, IsEnabled, Evaluate, EvaluateTrace, List)
5. **FeatureResult**: 판정 결과 + trace
6. **Gin 미들웨어**: Require, Inject
7. **테스트**: rule 함수 단위 테스트, Flags 통합 테스트, hashPercentage 분포 테스트

### 제외

- 피처 플래그 퍼시스턴스 (런타임 등록만, DB/파일 저장은 사용자 책임)
- 실시간 변경 (핫 리로드) — 별도 Phase에서 검토
- A/B 테스트 통계/분석
- 대시보드 UI

## 산출물

```
pkg/
  feature/
    user_context.go                — UserContext 구조체
    rule_is_beta_user.go           — IsBetaUser
    rule_is_internal_staff.go      — IsInternalStaff
    rule_is_region.go              — IsRegion (클로저)
    rule_has_attribute.go          — HasAttribute (클로저)
    rule_is_legacy_browser.go      — IsLegacyBrowser
    rule_is_user_in_percentage.go  — IsUserInPercentage (클로저)
    hash_percentage.go             — hashPercentage 결정론적 해시
    flags.go                       — Flags (NewFlags, Register, IsEnabled, Evaluate, EvaluateTrace, List)
    feature_result.go              — FeatureResult 구조체
    gin_feature.go                 — Require, Inject (Gin 미들웨어)
    rule_test.go                   — rule 함수 단위 테스트
    flags_test.go                  — Flags 통합 테스트
    hash_percentage_test.go        — 해시 분포 균등성 테스트
    gin_feature_test.go            — Gin 미들웨어 테스트
```

## 단계

### Step 1: 구조체 정의

- UserContext, FeatureResult

### Step 2: hashPercentage 구현

- 사용자 ID를 FNV-1a 등으로 해시하여 0.0~1.0 매핑
- 같은 ID는 항상 같은 값 반환 (결정론적)
- 분포 균등성 테스트

### Step 3: rule 함수 구현

- 각 rule 함수를 파일 하나에 하나씩 구현 (filefunc 규칙 준수)
- 각 함수는 1-2 depth 유지
- 클로저 rule: IsRegion, HasAttribute, IsUserInPercentage

### Step 4: Flags 구현

- NewFlags: 빈 레지스트리 생성
- Register: 피처 이름 + graph 등록
- IsEnabled: Evaluate → verdict >= 0이면 true
- Evaluate: graph.Evaluate(featureName, ctx) → verdict
- EvaluateTrace: graph.EvaluateTrace → FeatureResult
- List: 등록된 전체 피처를 순회하여 활성 목록 반환

### Step 5: Gin 미들웨어 구현

- Require: IsEnabled 판정 → 비활성이면 404 반환
- Inject: List 결과를 gin.Context에 저장

### Step 6: 테스트

- rule 함수 단위 테스트: 각 조건별 true/false
- hashPercentage 분포 테스트: 10000개 ID로 균등 분포 확인
- Flags 통합 테스트:
  - 베타 사용자 → 활성화
  - 레거시 브라우저 → 비활성화
  - 내부 직원 + 레거시 + defeat → 활성화
  - 30% 롤아웃 → 해시 기반 결정론적 결과
  - 미등록 피처 → 에러
  - List → 활성 피처만 반환
- Gin 미들웨어 테스트: Require 통과/거부, Inject 컨텍스트 주입

### Step 7: 전체 테스트 PASS 확인

- `go test ./...` 전체 PASS 확인

## 검증 기준

1. IsBetaUser, IsRegion 등 rule 함수가 ground에서 올바르게 판정한다
2. hashPercentage가 같은 ID에 대해 항상 같은 값을 반환한다 (결정론적)
3. hashPercentage 분포가 10000개 샘플에서 ±5% 이내로 균등하다
4. Flags.IsEnabled이 verdict >= 0이면 true를 반환한다
5. defeat edge가 예외를 정확히 처리한다 (내부 직원 → 레거시 제외 무시)
6. EvaluateTrace 결과에 각 rule의 판정 근거가 포함된다
7. Gin Require 미들웨어가 비활성 피처에 대해 404를 반환한다
8. 전체 테스트 PASS

## 의존성

- Phase 001-009: toulmin 코어 (NewGraph, Evaluate, EvaluateTrace)
