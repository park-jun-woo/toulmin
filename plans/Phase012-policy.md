# Phase 012: 정책 프레임워크 — pkg/policy

## 목표

toulmin 기반 정책 판정 프레임워크를 `pkg/policy`에 구현한다.
OPA/Rego 없이 Go 함수 + defeats graph로 인증/인가/접근 제어를 선언적으로 처리한다.
Gin 어댑터를 제공하여 미들웨어로 즉시 사용 가능하게 한다.
Phase 010의 backing 일급 시민 원칙을 적용하여, 판정 기준을 명시적 backing으로 선언한다.

## 배경

### 현재 문제

1. **OPA는 무겁다**: Rego 언어 학습 + OPA 런타임 + 번들 관리가 필요하다. 단순 정책에도 인프라 오버헤드가 크다
2. **미들웨어 순서 의존성**: Gin 미들웨어 체인에서 인증 → 인가 → rate limit → IP 차단 순서가 암묵적이다. 순서 변경 시 의도치 않은 동작이 발생한다
3. **예외 처리가 if-else**: 화이트리스트, 관리자 오버라이드 같은 예외가 미들웨어나 핸들러 안에 if-else로 박힌다. 규칙이 늘어나면 중첩이 깊어진다

### toulmin이 해결하는 것

- 정책 규칙 = Go 함수 (1-2 depth)
- 규칙 간 관계 = defeats graph (명시적 선언)
- 예외 = defeat edge (if-else 대체)
- 판정 근거 = EvaluateTrace (감사 로그, backing 값까지 추적)

### claim/ground/backing 3-element 분리 원칙

toulmin의 `func(claim any, ground any, backing any) (bool, any)` 시그니처가 프레임워크 확장성의 핵심이다.

- **claim = 뭘 판정하나**: 정책 프레임워크에서 claim은 요청 자체 (경로, 메서드, 리소스)
- **ground = 판정 재료**: ground는 판정에 필요한 런타임 컨텍스트 (사용자, IP 등). 평가 시점에 결정된다
- **backing = 판정 기준**: backing은 규칙의 판정 기준 (역할명, 화이트리스트, 임계값 등). 선언 시점에 고정된다

프레임워크는 graph 구조와 판정 흐름을 제공하고, **도메인 데이터는 ground로, 판정 기준은 backing으로 사용자가 주입한다.** rule 함수는 ground에서 사실을, backing에서 기준을 받아 판단만 하므로, 프레임워크가 도메인을 몰라도 동작한다.

| 역할 | 정책 프레임워크에서 |
|---|---|
| claim | 요청 (경로, 메서드, 리소스 ID) |
| ground | RequestContext (User, ClientIP) — 평가 시점 |
| backing | 판정 기준 (role명, IP 목록, 임계값) — 선언 시점 |
| rule 함수 | ground에서 사실, backing에서 기준을 받아 판단 (1-2 depth) |
| graph | rule 간 관계 선언 (defeat = 예외) |
| verdict | 허용/거부 판정 |

## 핵심 설계

### 범용 rule 함수 세트

모든 rule 함수는 `func(claim any, ground any, backing any) (bool, any)` 시그니처를 따른다. 클로저를 사용하지 않으며, 판정 기준은 backing으로 받는다.

```go
// pkg/policy/rule_is_authenticated.go
// backing: nil
func IsAuthenticated(claim any, ground any, backing any) (bool, any) {
    ctx := ground.(*RequestContext)
    return ctx.User != nil, nil
}

// pkg/policy/rule_is_in_role.go
// backing: string (role명)
func IsInRole(claim any, ground any, backing any) (bool, any) {
    ctx := ground.(*RequestContext)
    role := backing.(string)
    return ctx.User.Role == role, nil
}

// pkg/policy/rule_is_owner.go
// backing: nil
func IsOwner(claim any, ground any, backing any) (bool, any) {
    ctx := ground.(*RequestContext)
    return ctx.User.ID == ctx.ResourceOwnerID, nil
}

// pkg/policy/rule_is_ip_in_list.go
// backing: map[string]bool (IP 목록)
func IsIPInList(claim any, ground any, backing any) (bool, any) {
    ctx := ground.(*RequestContext)
    list := backing.(map[string]bool)
    return list[ctx.ClientIP], nil
}

// pkg/policy/rule_is_rate_limited.go
// backing: nil
func IsRateLimited(claim any, ground any, backing any) (bool, any) {
    ctx := ground.(*RequestContext)
    return ctx.RateLimiter.IsLimited(ctx.ClientIP), nil
}

// pkg/policy/rule_is_admin_override.go
// backing: nil
func IsAdminOverride(claim any, ground any, backing any) (bool, any) {
    ctx := ground.(*RequestContext)
    return ctx.User != nil && ctx.User.Role == "admin", nil
}
```

### RequestContext

```go
// pkg/policy/request_context.go
type RequestContext struct {
    User            *User
    ClientIP        string
    ResourceOwnerID string
    RateLimiter     RateLimiter
}
```

blocking IP 목록과 화이트리스트는 더 이상 RequestContext에 포함되지 않는다. 이들은 규칙의 판정 기준이므로 backing으로 선언 시점에 전달된다.

### 정책 선언

backing은 Warrant/Rebuttal/Defeater의 인자로 직접 전달한다. backing이 없는 rule은 nil을 명시한다.

```go
blockedIPs := map[string]bool{"1.2.3.4": true}
whitelistedIPs := map[string]bool{"10.0.0.1": true}

p := toulmin.NewGraph("admin:users").
    Warrant(policy.IsAuthenticated, nil, 1.0).
    Warrant(policy.IsInRole, "admin", 1.0).
    Rebuttal(policy.IsIPInList, blockedIPs, 1.0).
    Rebuttal(policy.IsRateLimited, nil, 1.0).
    Defeater(policy.IsIPInList, whitelistedIPs, 1.0).
    Defeater(policy.IsAdminOverride, nil, 1.0).
    DefeatWith(policy.IsIPInList, blockedIPs, policy.IsAuthenticated, nil).
    Defeat(policy.IsRateLimited, policy.IsAuthenticated).
    DefeatWith(policy.IsIPInList, whitelistedIPs, policy.IsIPInList, blockedIPs).
    Defeat(policy.IsAdminOverride, policy.IsRateLimited)
```

같은 `IsIPInList` 함수가 blockedIPs/whitelistedIPs 두 backing으로 사용된다. `DefeatWith`로 backing을 명시하여 구분한다. backing이 nil인 rule 간에는 `Defeat`를 사용한다.

### Gin 어댑터

Guard는 claim에 nil을 전달한다. 정책 판정에서 claim은 HTTP 요청 자체이며, 이미 라우트 매칭으로 확정되었기 때문이다. rule 함수가 필요로 하는 판정 재료는 ground(RequestContext)에, 판정 기준은 backing에 담긴다.

verdict `<= 0`이면 거부한다 (보안 컨텍스트: undecided는 거부).

```go
// pkg/policy/gin_guard.go
func Guard(g *toulmin.Graph, ctxBuilder ContextBuilderFunc) gin.HandlerFunc {
    return func(c *gin.Context) {
        ctx := ctxBuilder(c)
        results, err := g.EvaluateTrace(nil, ctx)
        if err != nil {
            c.AbortWithStatusJSON(500, gin.H{"error": "policy evaluation failed"})
            return
        }
        if results[0].Verdict <= 0 {
            c.AbortWithStatusJSON(403, gin.H{
                "error": "forbidden",
                "trace": formatTrace(results[0].Trace),
            })
            return
        }
        c.Next()
    }
}

// ContextBuilderFunc — gin.Context → RequestContext 변환
type ContextBuilderFunc func(*gin.Context) *RequestContext
```

### 사용 예시

```go
r := gin.Default()

blockedIPs := map[string]bool{"1.2.3.4": true}
whitelistedIPs := map[string]bool{"10.0.0.1": true}

adminPolicy := toulmin.NewGraph("admin:users").
    Warrant(policy.IsAuthenticated, nil, 1.0).
    Warrant(policy.IsInRole, "admin", 1.0).
    Rebuttal(policy.IsIPInList, blockedIPs, 1.0).
    Defeater(policy.IsIPInList, whitelistedIPs, 1.0).
    DefeatWith(policy.IsIPInList, blockedIPs, policy.IsAuthenticated, nil).
    DefeatWith(policy.IsIPInList, whitelistedIPs, policy.IsIPInList, blockedIPs)

r.GET("/admin/users", policy.Guard(adminPolicy, buildCtx), handler)
```

### 디버그 모드

```go
// 개발 환경에서 판정 근거를 응답 헤더로 노출
func GuardDebug(g *toulmin.Graph, ctxBuilder ContextBuilderFunc) gin.HandlerFunc
// X-Policy-Verdict: 0.33
// X-Policy-Trace: IsAuthenticated=true, IsInRole[admin]=true, IsIPInList[blockedIPs]=false
```

## 범위

### 포함

1. **RequestContext 구조체**: 정책 판정에 필요한 요청 컨텍스트 (런타임 데이터만, 판정 기준은 backing으로 분리)
2. **범용 rule 함수 세트**: IsAuthenticated, IsInRole, IsOwner, IsIPInList, IsRateLimited, IsAdminOverride
3. **Guard 함수**: toulmin graph + ContextBuilderFunc → gin.HandlerFunc 어댑터
4. **GuardDebug 함수**: 판정 근거를 헤더로 노출하는 디버그 버전
5. **formatTrace 헬퍼**: Trace 결과를 사람이 읽을 수 있는 문자열로 변환
6. **테스트**: rule 함수 단위 테스트, Guard 통합 테스트

### 제외

- RateLimiter 구현체 (인터페이스만 정의, 구현은 사용자)
- JWT 파싱 (ContextBuilderFunc에서 사용자가 처리)
- Echo/net/http 어댑터 — 별도 Phase에서 검토
- 정책 YAML 정의 및 코드젠 — fullend에서 처리

## 산출물

```
pkg/
  policy/
    request_context.go            — RequestContext 구조체 (런타임 데이터만)
    rate_limiter.go               — RateLimiter 인터페이스
    rule_is_authenticated.go      — IsAuthenticated (backing: nil)
    rule_is_in_role.go            — IsInRole (backing: string, role명)
    rule_is_owner.go              — IsOwner (backing: nil)
    rule_is_ip_in_list.go         — IsIPInList (backing: map[string]bool, IP 목록)
    rule_is_rate_limited.go       — IsRateLimited (backing: nil)
    rule_is_admin_override.go     — IsAdminOverride (backing: nil)
    gin_guard.go                  — Guard, GuardDebug
    format_trace.go               — formatTrace 헬퍼
    request_context_test.go       — 구조체 테스트
    rule_test.go                  — rule 함수 단위 테스트
    gin_guard_test.go             — Guard 통합 테스트
```

## 단계

### Step 1: RequestContext 및 인터페이스 정의

- `pkg/policy/request_context.go`: RequestContext, User 구조체 (런타임 데이터만 포함, IP 목록은 backing으로 분리)
- `pkg/policy/rate_limiter.go`: RateLimiter 인터페이스

### Step 2: rule 함수 구현

- 각 rule 함수를 파일 하나에 하나씩 구현 (filefunc 규칙 준수)
- 모든 함수는 `func(claim any, ground any, backing any) (bool, any)` 시그니처
- 각 함수는 1-2 depth 유지
- 판정 기준은 backing에서 type assertion으로 꺼낸다

### Step 3: Guard 어댑터 구현

- `pkg/policy/gin_guard.go`: Guard, GuardDebug
- `pkg/policy/format_trace.go`: trace 포매팅
- verdict `<= 0`이면 거부 (보안 컨텍스트)

### Step 4: 테스트

- rule 함수 단위 테스트: 각 rule이 올바른 claim, ground, backing 조합에서 true/false 반환
- Guard 통합 테스트: httptest로 Gin 라우터 생성, 정책 통과/거부 시나리오 검증
- 디버그 모드: 응답 헤더에 trace 포함 확인
- DefeatWith 시나리오: 같은 함수 + 다른 backing 간 defeat 동작 검증

### Step 5: 전체 테스트 PASS 확인

- `go test ./...` 전체 PASS 확인

## 검증 기준

1. IsAuthenticated, IsInRole, IsOwner 등 rule 함수가 `(claim, ground, backing)` 시그니처로 올바르게 판정한다
2. IsInRole이 backing으로 role명을 받아 재사용 가능하다 (클로저 없음)
3. IsIPInList이 backing으로 IP 목록을 받아 차단/화이트리스트를 구분한다
4. Guard가 verdict `<= 0`이면 403을 반환한다 (보안 컨텍스트: undecided는 거부)
5. GuardDebug가 판정 근거를 응답 헤더에 포함한다
6. DefeatWith로 같은 함수 + 다른 backing 간 defeat edge가 정확히 동작한다
7. Defeat는 backing nil인 rule 간에만 사용된다
8. EvaluateTrace 결과에 backing 값이 포함되어 감사 로그로 사용 가능하다
9. RequestContext에 판정 기준(IP 목록 등)이 포함되지 않는다 (backing으로 분리)
10. 전체 테스트 PASS

## 의존성

- Phase 001-009: toulmin 코어 (NewGraph, Evaluate, EvaluateTrace)
- Phase 010: backing 일급 시민 (3-element 시그니처, DefeatWith, funcID+backing 식별)
