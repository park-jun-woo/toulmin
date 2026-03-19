# Phase 011: 라우트 판정 프레임워크 — pkg/route (구현 완료)

## 목표

toulmin 기반 라우트 판정 프레임워크를 `pkg/route`에 구현한다.
Gin 미들웨어 체인의 순서 의존성을 defeats graph로 대체한다.
인증, 인가, rate limiting, IP 차단과 그 예외를 하나의 graph로 선언적 처리한다.

## 배경

### 현재 문제

1. **미들웨어 순서가 암묵적이다**: Gin에서 `AuthMiddleware(), RoleMiddleware(), RateLimitMiddleware()` 순서가 동작을 결정한다. 순서 변경 시 의도치 않은 동작이 발생한다
2. **핸들러 안에 조건 분기가 남는다**: 미들웨어를 통과해도 핸들러 안에서 `if isOwner || (isAdmin && !isSuspended)` 같은 분기가 반복된다
3. **예외 처리가 미들웨어 사이에 끼어든다**: 화이트리스트, 관리자 오버라이드, 내부 서비스 예외가 미들웨어 체인의 복잡도를 높인다

### toulmin이 해결하는 것

- 라우트 판정 조건 = rule 함수 (1-2 depth)
- 미들웨어 순서 의존성 → defeats graph (관계 선언, 순서 무관)
- 예외 = defeat edge (if-else 대체)
- 판정 근거 = EvaluateTrace ("왜 403인가" 즉시 추적)

### claim/ground/backing 분리 원칙

toulmin의 `(claim any, ground any, backing any)` 시그니처가 프레임워크 확장성의 핵심이다.

- **claim = 뭘 판정하나**: 라우트 프레임워크에서 claim은 nil (라우트 매칭으로 이미 확정)
- **ground = 판정 대상의 사실**: 요청 컨텍스트 (사용자, IP 등). 요청마다 달라진다
- **backing = 규칙의 판정 기준**: 역할명, 차단 목록, rate limiter 등. graph 선언 시 고정된다

| 역할 | 라우트 프레임워크에서 |
|---|---|
| claim | nil (라우트 매칭으로 확정) |
| ground | RouteContext (User, ClientIP, Headers, ...) |
| backing | 판정 기준 ("admin", blocklist, limiter, nil) |
| rule 함수 | ground에서 사실 확인, backing에서 기준 참조 (1-2 depth) |
| graph | rule 간 관계 선언 (defeat = 예외) |
| verdict | 요청 허용/거부 판정 (> 0 허용, <= 0 거부) |

## 핵심 설계

### RouteContext

```go
// pkg/route/route_context.go
type RouteContext struct {
    User     *User
    ClientIP string
    Method   string
    Path     string
    Headers  map[string]string
    Metadata map[string]any
}

// pkg/route/user.go
type User struct {
    ID    string
    Role  string
    Email string
}
```

### 범용 rule 함수 — backing 방식 (클로저 불필요)

Phase 010에서 backing이 일급 시민이 되면서, 판정 기준은 backing으로 전달한다. 클로저 팩토리 패턴이 사라지고 순수 함수만 남는다.

```go
// pkg/route/rule_is_authenticated.go — backing 불필요 (nil)
func IsAuthenticated(claim any, ground any, backing any) (bool, any) {
    ctx := ground.(*RouteContext)
    return ctx.User != nil, nil
}

// pkg/route/rule_has_role.go — backing: string (역할명)
func HasRole(claim any, ground any, backing any) (bool, any) {
    ctx := ground.(*RouteContext)
    role := backing.(string)
    return ctx.User != nil && ctx.User.Role == role, nil
}

// pkg/route/rule_is_owner.go — backing: func(*RouteContext) string (소유자 ID 추출)
func IsOwner(claim any, ground any, backing any) (bool, any) {
    ctx := ground.(*RouteContext)
    if ctx.User == nil {
        return false, nil
    }
    ownerIDFunc := backing.(func(*RouteContext) string)
    return ctx.User.ID == ownerIDFunc(ctx), nil
}

// pkg/route/rule_is_ip_blocked.go — backing: func(string) bool (차단 목록)
func IsIPBlocked(claim any, ground any, backing any) (bool, any) {
    ctx := ground.(*RouteContext)
    blocklist := backing.(func(string) bool)
    return blocklist(ctx.ClientIP), nil
}

// pkg/route/rule_is_whitelisted.go — backing: func(string) bool (화이트리스트)
func IsWhitelisted(claim any, ground any, backing any) (bool, any) {
    ctx := ground.(*RouteContext)
    whitelist := backing.(func(string) bool)
    return whitelist(ctx.ClientIP), nil
}

// pkg/route/rule_is_rate_limited.go — backing: RateLimiter
func IsRateLimited(claim any, ground any, backing any) (bool, any) {
    ctx := ground.(*RouteContext)
    limiter := backing.(RateLimiter)
    return limiter.IsLimited(ctx.ClientIP), nil
}

// pkg/route/rule_is_internal_service.go — backing 불필요 (nil)
func IsInternalService(claim any, ground any, backing any) (bool, any) {
    ctx := ground.(*RouteContext)
    token := ctx.Headers["X-Internal-Token"]
    return token != "", nil
}

// pkg/route/rule_is_admin_override.go — backing 불필요 (nil)
func IsAdminOverride(claim any, ground any, backing any) (bool, any) {
    ctx := ground.(*RouteContext)
    return ctx.User != nil && ctx.User.Role == "admin", nil
}
```

### RateLimiter 인터페이스

```go
// pkg/route/rate_limiter.go
type RateLimiter interface {
    IsLimited(key string) bool
}
```

### Guard — Gin 미들웨어

Guard는 claim에 nil을 전달한다. verdict <= 0이면 거부 (보안 컨텍스트: undecided는 거부).

```go
// pkg/route/guard.go
func Guard(g *toulmin.GraphBuilder, ctxBuilder ContextBuilderFunc) gin.HandlerFunc {
    return func(c *gin.Context) {
        ctx := ctxBuilder(c)
        results, err := g.EvaluateTrace(nil, ctx)
        if err != nil {
            c.AbortWithStatusJSON(500, gin.H{"error": "route evaluation failed"})
            return
        }
        if len(results) == 0 || results[0].Verdict <= 0 {
            c.AbortWithStatusJSON(403, gin.H{"error": "forbidden"})
            return
        }
        c.Next()
    }
}

// ContextBuilderFunc — gin.Context → RouteContext 변환
type ContextBuilderFunc func(*gin.Context) *RouteContext
```

### GuardDebug — 디버그 미들웨어

```go
// pkg/route/guard_debug.go
func GuardDebug(g *toulmin.GraphBuilder, ctxBuilder ContextBuilderFunc) gin.HandlerFunc
// 판정 근거를 응답 헤더로 노출
// X-Route-Verdict: 0.33
// X-Route-Trace: IsAuthenticated=true, HasRole=true, IsIPBlocked=false
```

### 사용 예시

backing이 일급 시민이므로 클로저가 불필요하다. 같은 함수를 다른 backing으로 재사용한다. backing이 있는 rule 간 Defeat는 `DefeatWith`를 사용한다.

```go
r := gin.Default()

blocklist := func(ip string) bool { return ip == "1.2.3.4" }
whitelist := func(ip string) bool { return ip == "10.0.0.1" }

adminRoute := toulmin.NewGraph("route:admin").
    Warrant(route.IsAuthenticated, nil, 1.0).
    Warrant(route.HasRole, "admin", 1.0).
    Rebuttal(route.IsIPBlocked, blocklist, 1.0).
    Rebuttal(route.IsRateLimited, limiter, 1.0).
    Defeater(route.IsWhitelisted, whitelist, 1.0).
    Defeater(route.IsInternalService, nil, 1.0).
    DefeatWith(route.IsIPBlocked, blocklist, route.IsAuthenticated, nil).
    DefeatWith(route.IsRateLimited, limiter, route.IsAuthenticated, nil).
    DefeatWith(route.IsWhitelisted, whitelist, route.IsIPBlocked, blocklist).
    DefeatWith(route.IsInternalService, nil, route.IsRateLimited, limiter)

publicRoute := toulmin.NewGraph("route:public").
    Warrant(route.IsAuthenticated, nil, 1.0).
    Rebuttal(route.IsIPBlocked, blocklist, 1.0).
    DefeatWith(route.IsIPBlocked, blocklist, route.IsAuthenticated, nil)

r.GET("/admin/users", route.Guard(adminRoute, buildCtx), adminHandler)
r.GET("/posts", route.Guard(publicRoute, buildCtx), postsHandler)
```

## 범위

### 포함

1. **RouteContext, User 구조체**: 라우트 판정에 필요한 요청 컨텍스트
2. **RateLimiter 인터페이스**: rate limiting 추상화
3. **범용 rule 함수**: IsAuthenticated, HasRole, IsOwner, IsIPBlocked, IsWhitelisted, IsRateLimited, IsInternalService, IsAdminOverride — 전부 backing 방식, 클로저 없음
4. **Guard**: toulmin graph → gin.HandlerFunc 어댑터 (verdict <= 0 거부)
5. **GuardDebug**: 판정 근거를 헤더로 노출하는 디버그 버전
6. **테스트**: rule 함수 단위 테스트, Guard 통합 테스트

### 제외

- RateLimiter 구현체 (인터페이스만 정의, 구현은 사용자)
- JWT 파싱 (ContextBuilderFunc에서 사용자가 처리)
- Echo/net/http 어댑터 — 별도 Phase에서 검토
- CORS, 로깅 등 판정과 무관한 미들웨어

## 산출물

```
pkg/
  route/
    route_context.go              — RouteContext 구조체
    user.go                       — User 구조체
    rate_limiter.go               — RateLimiter 인터페이스
    rule_is_authenticated.go      — IsAuthenticated (backing: nil)
    rule_has_role.go              — HasRole (backing: string)
    rule_is_owner.go              — IsOwner (backing: func)
    rule_is_ip_blocked.go         — IsIPBlocked (backing: func)
    rule_is_whitelisted.go        — IsWhitelisted (backing: func)
    rule_is_rate_limited.go       — IsRateLimited (backing: RateLimiter)
    rule_is_internal_service.go   — IsInternalService (backing: nil)
    rule_is_admin_override.go     — IsAdminOverride (backing: nil)
    guard.go                      — Guard (Gin 미들웨어, verdict <= 0 거부)
    guard_debug.go                — GuardDebug (디버그 미들웨어)
    rule_test.go                  — rule 함수 단위 테스트
    guard_test.go                 — Guard 통합 테스트
```

## 구현 과정 발견 사항

Phase 010(backing) 이전에 먼저 구현되어 클로저 방식으로 작성 후, backing 도입 시 순수 함수 방식으로 전환됨. 구현 과정에서 발견된 핵심 패턴:

1. **Rebuttal만으로는 공격 안 됨** — 반드시 Defeat/DefeatWith edge 선언 필요
2. **Defeater 등록 필수** — Defeat에서 참조하는 함수는 반드시 Warrant/Rebuttal/Defeater로 등록
3. **verdict 0.0 = undecided** — 보안 컨텍스트에서 <= 0 거부가 안전한 기본값
4. **backing이 있는 rule 간 Defeat** — 같은 함수 + 다른 backing은 DefeatWith로 명시

## 검증 기준

1. IsAuthenticated, HasRole, IsOwner 등 rule 함수가 backing을 통해 올바르게 판정한다
2. HasRole은 backing으로 역할명을 받아 클로저 없이 재사용 가능하다
3. Guard가 verdict > 0이면 c.Next(), <= 0이면 403을 반환한다 (보안 컨텍스트: undecided는 거부)
4. GuardDebug가 판정 근거를 응답 헤더에 포함한다
5. DefeatWith가 backing이 있는 rule 간 defeat를 정확히 처리한다
6. 미들웨어 순서에 관계없이 graph 선언대로 판정된다
7. 전체 테스트 PASS

## 의존성

- Phase 001-009: toulmin 코어 (NewGraph, Evaluate, EvaluateTrace)
- Phase 010: backing 일급 시민 (Warrant(fn, backing, qualifier), DefeatWith)
