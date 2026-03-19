# Phase 010: 정책 프레임워크 — pkg/policy

## 목표

toulmin 기반 정책 판정 프레임워크를 `pkg/policy`에 구현한다.
OPA/Rego 없이 Go 함수 + defeats graph로 인증/인가/접근 제어를 선언적으로 처리한다.
Gin 어댑터를 제공하여 미들웨어로 즉시 사용 가능하게 한다.

## 배경

### 현재 문제

1. **OPA는 무겁다**: Rego 언어 학습 + OPA 런타임 + 번들 관리가 필요하다. 단순 정책에도 인프라 오버헤드가 크다
2. **미들웨어 순서 의존성**: Gin 미들웨어 체인에서 인증 → 인가 → rate limit → IP 차단 순서가 암묵적이다. 순서 변경 시 의도치 않은 동작이 발생한다
3. **예외 처리가 if-else**: 화이트리스트, 관리자 오버라이드 같은 예외가 미들웨어나 핸들러 안에 if-else로 박힌다. 규칙이 늘어나면 중첩이 깊어진다

### toulmin이 해결하는 것

- 정책 규칙 = Go 함수 (1-2 depth)
- 규칙 간 관계 = defeats graph (명시적 선언)
- 예외 = defeat edge (if-else 대체)
- 판정 근거 = EvaluateTrace (감사 로그)

### claim/ground 분리 원칙

toulmin의 `(claim any, ground any)` 시그니처가 프레임워크 확장성의 핵심이다.

- **claim = 뭘 판정하나**: 정책 프레임워크에서 claim은 요청 자체 (경로, 메서드, 리소스)
- **ground = 판정 재료**: ground는 판정에 필요한 컨텍스트 (사용자, IP, 역할 등)

프레임워크는 graph 구조와 판정 흐름을 제공하고, **도메인 데이터는 ground로 사용자가 주입한다.** rule 함수는 ground에서 데이터를 꺼내 판단만 하므로, 프레임워크가 도메인을 몰라도 동작한다.

| 역할 | 정책 프레임워크에서 |
|---|---|
| claim | 요청 (경로, 메서드, 리소스 ID) |
| ground | RequestContext (User, ClientIP, BlockedIPs, ...) |
| rule 함수 | ground에서 조건 하나만 판단 (1-2 depth) |
| graph | rule 간 관계 선언 (defeat = 예외) |
| verdict | 허용/거부 판정 |

## 핵심 설계

### 범용 rule 함수 세트

```go
// pkg/policy/rule_is_authenticated.go
func IsAuthenticated(claim any, ground any) (bool, any) {
    ctx := ground.(*RequestContext)
    return ctx.User != nil, nil
}

// pkg/policy/rule_has_role.go
func HasRole(role string) toulmin.RuleFunc {
    return func(claim any, ground any) (bool, any) {
        ctx := ground.(*RequestContext)
        return ctx.User.Role == role, nil
    }
}

// pkg/policy/rule_is_owner.go
func IsOwner(claim any, ground any) (bool, any) {
    ctx := ground.(*RequestContext)
    return ctx.User.ID == ctx.ResourceOwnerID, nil
}

// pkg/policy/rule_is_ip_blocked.go
func IsIPBlocked(claim any, ground any) (bool, any) {
    ctx := ground.(*RequestContext)
    return ctx.BlockedIPs[ctx.ClientIP], nil
}

// pkg/policy/rule_is_rate_limited.go
func IsRateLimited(claim any, ground any) (bool, any) {
    ctx := ground.(*RequestContext)
    return ctx.RateLimiter.IsLimited(ctx.ClientIP), nil
}

// pkg/policy/rule_is_whitelisted.go
func IsWhitelisted(claim any, ground any) (bool, any) {
    ctx := ground.(*RequestContext)
    return ctx.WhitelistedIPs[ctx.ClientIP], nil
}

// pkg/policy/rule_is_admin_override.go
func IsAdminOverride(claim any, ground any) (bool, any) {
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
    BlockedIPs      map[string]bool
    WhitelistedIPs  map[string]bool
    RateLimiter     RateLimiter
}
```

### 정책 선언

**주의**: 클로저 rule은 변수에 저장 후 재사용해야 한다. Rebuttal만으로는 공격이 일어나지 않으며 반드시 Defeat edge를 선언해야 한다. 예외를 처리하는 rule은 Defeater로 등록해야 한다.

```go
// 클로저는 변수에 저장 후 재사용
hasAdmin := policy.HasRole("admin")

p := toulmin.NewGraph("admin:users").
    Warrant(policy.IsAuthenticated, 1.0).
    Warrant(hasAdmin, 1.0).
    Rebuttal(policy.IsIPBlocked, 1.0).
    Rebuttal(policy.IsRateLimited, 1.0).
    Defeater(policy.IsWhitelisted, 1.0).       // 예외 rule은 Defeater로 등록
    Defeater(policy.IsAdminOverride, 1.0).
    Defeat(policy.IsIPBlocked, policy.IsAuthenticated).  // Rebuttal → Warrant 공격 edge 필수
    Defeat(policy.IsRateLimited, policy.IsAuthenticated).
    Defeat(policy.IsWhitelisted, policy.IsIPBlocked).    // Defeater → Rebuttal 예외 처리
    Defeat(policy.IsAdminOverride, policy.IsRateLimited)
```

### Gin 어댑터

Guard는 claim에 nil을 전달한다. 정책 판정에서 claim은 HTTP 요청 자체이며, 이미 라우트 매칭으로 확정되었기 때문이다. rule 함수가 필요로 하는 모든 판정 재료는 ground(RequestContext)에 담긴다.

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

hasAdmin := policy.HasRole("admin")

adminPolicy := toulmin.NewGraph("admin:users").
    Warrant(policy.IsAuthenticated, 1.0).
    Warrant(hasAdmin, 1.0).
    Rebuttal(policy.IsIPBlocked, 1.0).
    Defeater(policy.IsWhitelisted, 1.0).
    Defeat(policy.IsIPBlocked, policy.IsAuthenticated).
    Defeat(policy.IsWhitelisted, policy.IsIPBlocked)

r.GET("/admin/users", policy.Guard(adminPolicy, buildCtx), handler)
```

### 디버그 모드

```go
// 개발 환경에서 판정 근거를 응답 헤더로 노출
func GuardDebug(g *toulmin.Graph, ctxBuilder ContextBuilderFunc) gin.HandlerFunc
// X-Policy-Verdict: 0.33
// X-Policy-Trace: IsAuthenticated=true, HasRole=true, IsIPBlocked=false
```

## 범위

### 포함

1. **RequestContext 구조체**: 정책 판정에 필요한 요청 컨텍스트
2. **범용 rule 함수 세트**: IsAuthenticated, HasRole, IsOwner, IsIPBlocked, IsRateLimited, IsWhitelisted, IsAdminOverride
3. **Guard 함수**: toulmin graph → gin.HandlerFunc 어댑터
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
    request_context.go            — RequestContext 구조체
    rate_limiter.go               — RateLimiter 인터페이스
    rule_is_authenticated.go      — IsAuthenticated
    rule_has_role.go              — HasRole (클로저)
    rule_is_owner.go              — IsOwner
    rule_is_ip_blocked.go         — IsIPBlocked
    rule_is_rate_limited.go       — IsRateLimited
    rule_is_whitelisted.go        — IsWhitelisted
    rule_is_admin_override.go     — IsAdminOverride
    gin_guard.go                  — Guard, GuardDebug
    format_trace.go               — formatTrace 헬퍼
    request_context_test.go       — 구조체 테스트
    rule_test.go                  — rule 함수 단위 테스트
    gin_guard_test.go             — Guard 통합 테스트
```

## 단계

### Step 1: RequestContext 및 인터페이스 정의

- `pkg/policy/request_context.go`: RequestContext, User 구조체
- `pkg/policy/rate_limiter.go`: RateLimiter 인터페이스

### Step 2: rule 함수 구현

- 각 rule 함수를 파일 하나에 하나씩 구현 (filefunc 규칙 준수)
- 각 함수는 1-2 depth 유지

### Step 3: Guard 어댑터 구현

- `pkg/policy/gin_guard.go`: Guard, GuardDebug
- `pkg/policy/format_trace.go`: trace 포매팅

### Step 4: 테스트

- rule 함수 단위 테스트: 각 rule이 올바른 조건에서 true/false 반환
- Guard 통합 테스트: httptest로 Gin 라우터 생성, 정책 통과/거부 시나리오 검증
- 디버그 모드: 응답 헤더에 trace 포함 확인

### Step 5: 전체 테스트 PASS 확인

- `go test ./...` 전체 PASS 확인

## 검증 기준

1. IsAuthenticated, HasRole, IsOwner 등 rule 함수가 올바르게 판정한다
2. HasRole이 클로저로 역할명을 받아 재사용 가능하다
3. Guard가 verdict > 0이면 c.Next(), <= 0이면 403을 반환한다 (보안 컨텍스트: undecided는 거부)
4. GuardDebug가 판정 근거를 응답 헤더에 포함한다
5. defeat edge가 예외를 정확히 처리한다 (화이트리스트 → IP 차단 무시)
6. EvaluateTrace 결과가 감사 로그로 사용 가능하다
7. 전체 테스트 PASS

## 의존성

- Phase 001-009: toulmin 코어 (NewGraph, Evaluate, EvaluateTrace)
