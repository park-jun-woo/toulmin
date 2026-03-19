# Phase 010: 라우트 판정 프레임워크 — pkg/route

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

### claim/ground 분리 원칙

toulmin의 `(claim any, ground any)` 시그니처가 프레임워크 확장성의 핵심이다.

- **claim = 뭘 판정하나**: 라우트 프레임워크에서 claim은 nil (라우트 매칭으로 이미 확정)
- **ground = 판정 재료**: ground는 요청 컨텍스트 (사용자, IP, 역할 등)

프레임워크는 Guard 함수와 판정 흐름을 제공하고, **도메인 데이터는 ground로 사용자가 주입한다.** rule 함수는 ground에서 데이터를 꺼내 판단만 하므로, 프레임워크가 도메인을 몰라도 동작한다.

| 역할 | 라우트 프레임워크에서 |
|---|---|
| claim | nil (라우트 매칭으로 확정) |
| ground | RouteContext (User, ClientIP, Headers, ...) |
| rule 함수 | ground에서 조건 하나만 판단 (1-2 depth) |
| graph | rule 간 관계 선언 (defeat = 예외) |
| verdict | 요청 허용/거부 판정 |

라우트 프레임워크는 가장 기본적인 판정 레이어다. 이후 Phase의 policy, state, approve 등이 이 위에 쌓인다. claim이 nil인 것은 라우트 판정의 특성 — HTTP 요청이 어떤 엔드포인트에 도달했는지는 Gin 라우터가 이미 결정했으므로, "이 요청을 통과시킬 것인가"만 판정하면 된다.

## 핵심 설계

### RouteContext

```go
// pkg/route/route_context.go
type RouteContext struct {
    User       *User
    ClientIP   string
    Method     string
    Path       string
    Headers    map[string]string
    Metadata   map[string]any
}

// pkg/route/user.go
type User struct {
    ID    string
    Role  string
    Email string
}
```

### 범용 rule 함수

```go
// pkg/route/rule_is_authenticated.go
func IsAuthenticated(claim any, ground any) (bool, any) {
    ctx := ground.(*RouteContext)
    return ctx.User != nil, nil
}

// pkg/route/rule_has_role.go
func HasRole(role string) toulmin.RuleFunc {
    return func(claim any, ground any) (bool, any) {
        ctx := ground.(*RouteContext)
        return ctx.User != nil && ctx.User.Role == role, nil
    }
}

// pkg/route/rule_is_owner.go
func IsOwner(ownerIDFunc func(*RouteContext) string) toulmin.RuleFunc {
    return func(claim any, ground any) (bool, any) {
        ctx := ground.(*RouteContext)
        if ctx.User == nil {
            return false, nil
        }
        return ctx.User.ID == ownerIDFunc(ctx), nil
    }
}

// pkg/route/rule_is_ip_blocked.go
func IsIPBlocked(blocklist func(string) bool) toulmin.RuleFunc {
    return func(claim any, ground any) (bool, any) {
        ctx := ground.(*RouteContext)
        return blocklist(ctx.ClientIP), nil
    }
}

// pkg/route/rule_is_whitelisted.go
func IsWhitelisted(whitelist func(string) bool) toulmin.RuleFunc {
    return func(claim any, ground any) (bool, any) {
        ctx := ground.(*RouteContext)
        return whitelist(ctx.ClientIP), nil
    }
}

// pkg/route/rule_is_rate_limited.go
func IsRateLimited(limiter RateLimiter) toulmin.RuleFunc {
    return func(claim any, ground any) (bool, any) {
        ctx := ground.(*RouteContext)
        return limiter.IsLimited(ctx.ClientIP), nil
    }
}

// pkg/route/rule_is_internal_service.go
func IsInternalService(claim any, ground any) (bool, any) {
    ctx := ground.(*RouteContext)
    token, _ := ctx.Headers["X-Internal-Token"]
    return token != "", nil
}

// pkg/route/rule_is_admin_override.go
func IsAdminOverride(claim any, ground any) (bool, any) {
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

```go
// pkg/route/guard.go
func Guard(g *toulmin.Graph, ctxBuilder ContextBuilderFunc) gin.HandlerFunc {
    return func(c *gin.Context) {
        ctx := ctxBuilder(c)
        results, err := g.EvaluateTrace(nil, ctx)
        if err != nil {
            c.AbortWithStatusJSON(500, gin.H{"error": "route evaluation failed"})
            return
        }
        if results[0].Verdict < 0 {
            c.AbortWithStatusJSON(403, gin.H{
                "error": "forbidden",
            })
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
func GuardDebug(g *toulmin.Graph, ctxBuilder ContextBuilderFunc) gin.HandlerFunc
// 판정 근거를 응답 헤더로 노출
// X-Route-Verdict: 0.33
// X-Route-Trace: IsAuthenticated=true, HasRole=true, IsIPBlocked=false
```

### 사용 예시

**주의**: 클로저 rule은 변수에 저장 후 재사용해야 한다. 호출할 때마다 새 인스턴스가 생성되어 funcID가 달라지므로 Rebuttal/Defeat 간 매칭이 안 된다. Rebuttal만으로는 공격이 일어나지 않으며 반드시 Defeat edge를 선언해야 한다. 예외를 처리하는 rule은 Defeater로 등록해야 한다.

```go
r := gin.Default()

blocklist := func(ip string) bool { return ip == "1.2.3.4" }
whitelist := func(ip string) bool { return ip == "10.0.0.1" }

// 클로저는 변수에 저장 후 재사용
hasAdmin := route.HasRole("admin")
ipBlocked := route.IsIPBlocked(blocklist)
whitelisted := route.IsWhitelisted(whitelist)
rateLimited := route.IsRateLimited(limiter)

adminRoute := toulmin.NewGraph("route:admin").
    Warrant(route.IsAuthenticated, 1.0).
    Warrant(hasAdmin, 1.0).
    Rebuttal(ipBlocked, 1.0).
    Rebuttal(rateLimited, 1.0).
    Defeater(whitelisted, 1.0).              // 예외 rule은 Defeater로 등록
    Defeater(route.IsInternalService, 1.0).
    Defeat(ipBlocked, route.IsAuthenticated). // Rebuttal → Warrant 공격 edge 필수
    Defeat(rateLimited, route.IsAuthenticated).
    Defeat(whitelisted, ipBlocked).           // Defeater → Rebuttal 예외 처리
    Defeat(route.IsInternalService, rateLimited)

publicRoute := toulmin.NewGraph("route:public").
    Warrant(route.IsAuthenticated, 1.0).
    Rebuttal(ipBlocked, 1.0).
    Defeat(ipBlocked, route.IsAuthenticated)

r.GET("/admin/users", route.Guard(adminRoute, buildCtx), adminHandler)
r.GET("/posts", route.Guard(publicRoute, buildCtx), postsHandler)
```

## 범위

### 포함

1. **RouteContext, User 구조체**: 라우트 판정에 필요한 요청 컨텍스트
2. **RateLimiter 인터페이스**: rate limiting 추상화
3. **범용 rule 함수**: IsAuthenticated, HasRole, IsOwner, IsIPBlocked, IsWhitelisted, IsRateLimited, IsInternalService, IsAdminOverride
4. **Guard**: toulmin graph → gin.HandlerFunc 어댑터
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
    rule_is_authenticated.go      — IsAuthenticated
    rule_has_role.go              — HasRole (클로저)
    rule_is_owner.go              — IsOwner (클로저)
    rule_is_ip_blocked.go         — IsIPBlocked (클로저)
    rule_is_whitelisted.go        — IsWhitelisted (클로저)
    rule_is_rate_limited.go       — IsRateLimited (클로저)
    rule_is_internal_service.go   — IsInternalService
    rule_is_admin_override.go     — IsAdminOverride
    guard.go                      — Guard (Gin 미들웨어)
    guard_debug.go                — GuardDebug (디버그 미들웨어)
    rule_test.go                  — rule 함수 단위 테스트
    guard_test.go                 — Guard 통합 테스트
```

## 단계

### Step 1: 구조체 및 인터페이스 정의

- `pkg/route/route_context.go`: RouteContext
- `pkg/route/user.go`: User
- `pkg/route/rate_limiter.go`: RateLimiter 인터페이스

### Step 2: rule 함수 구현

- 각 rule 함수를 파일 하나에 하나씩 구현 (filefunc 규칙 준수)
- 각 함수는 1-2 depth 유지
- 클로저 rule: HasRole, IsOwner, IsIPBlocked, IsWhitelisted, IsRateLimited

### Step 3: Guard 어댑터 구현

- `pkg/route/guard.go`: Guard — graph.EvaluateTrace(nil, ctx) → verdict 판정 → c.Next() 또는 c.Abort
- `pkg/route/guard_debug.go`: GuardDebug — 판정 근거를 X-Route-Verdict, X-Route-Trace 헤더로 노출

### Step 4: 테스트

- rule 함수 단위 테스트: 각 rule이 올바른 조건에서 true/false 반환
- Guard 통합 테스트: httptest로 Gin 라우터 생성
  - 인증된 admin → 200
  - 미인증 → 403
  - IP 차단 → 403
  - IP 차단 + 화이트리스트 defeat → 200
  - rate limited → 403
  - rate limited + 내부 서비스 defeat → 200
- GuardDebug 테스트: 응답 헤더에 verdict, trace 포함 확인

### Step 5: 전체 테스트 PASS 확인

- `go test ./...` 전체 PASS 확인

## 검증 기준

1. IsAuthenticated, HasRole, IsOwner 등 rule 함수가 올바르게 판정한다
2. HasRole이 클로저로 역할명을 받아 재사용 가능하다
3. Guard가 verdict > 0이면 c.Next(), <= 0이면 403을 반환한다 (보안 컨텍스트: undecided는 거부)
4. GuardDebug가 판정 근거를 응답 헤더에 포함한다
5. defeat edge가 예외를 정확히 처리한다 (화이트리스트 → IP 차단 무시, 내부 서비스 → rate limit 무시)
6. 미들웨어 순서에 관계없이 graph 선언대로 판정된다
7. 전체 테스트 PASS

## 의존성

- Phase 001-009: toulmin 코어 (NewGraph, Evaluate, EvaluateTrace)
