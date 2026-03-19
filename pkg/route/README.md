# pkg/route

**Stop chaining middleware. Declare rules and relationships.**

HTTP route guard built on toulmin defeats graph. Replaces middleware ordering with declarative rule relationships. Works with Gin.

## Install

```go
import "github.com/park-jun-woo/toulmin/pkg/route"
```

## Quick Start

```go
blocklist := func(ip string) bool { return ip == "1.2.3.4" }
whitelist := func(ip string) bool { return ip == "10.0.0.1" }

g := toulmin.NewGraph("route:admin")
auth := g.Warrant(route.IsAuthenticated, nil, 1.0)
admin := g.Warrant(route.IsInRole, "admin", 1.0)
blocked := g.Rebuttal(route.IsIPInList, blocklist, 1.0)
allowed := g.Defeater(route.IsIPInList, whitelist, 1.0)
g.Defeat(blocked, auth)
g.Defeat(allowed, blocked)

r := gin.Default()
r.GET("/admin/users", route.Guard(g, buildCtx), handler)
```

## Rules

All rules follow the Toulmin signature: `func(claim, ground, backing) (bool, any)`.

- **ground** = request facts (User, IP, headers) — changes per request
- **backing** = judgment criteria (role name, IP list, rate limiter) — fixed at declaration

| Rule | Backing | Description |
|---|---|---|
| `IsAuthenticated` | nil | User is not nil |
| `IsInRole` | string | User role matches backing |
| `IsOwner` | func(*RouteContext) string | User ID matches resource owner |
| `IsIPInList` | func(string) bool | Client IP in list (blocklist or whitelist) |
| `IsRateLimited` | RateLimiter | Client IP is rate limited |
| `IsInternalService` | nil | X-Internal-Token header exists |
| `IsAdminOverride` | nil | User role is admin |

## Guard vs GuardDebug

| | Guard | GuardDebug |
|---|---|---|
| Evaluation | `Evaluate` (lightweight) | `EvaluateTrace` (with trace) |
| Response headers | None | `X-Route-Verdict`, `X-Route-Trace` |
| Use case | Production | Development/debugging |

Both deny when `verdict <= 0` (undecided = deny in security context).

## Same Function, Different Backing

`IsIPInList` serves both blocklist and whitelist — same function, different backing. Use `*Rule` references to wire defeat relationships:

```go
g := toulmin.NewGraph("firewall")
auth := g.Warrant(route.IsAuthenticated, nil, 1.0)
blocked := g.Rebuttal(route.IsIPInList, blocklist, 1.0)
allowed := g.Defeater(route.IsIPInList, whitelist, 1.0)
g.Defeat(blocked, auth)
g.Defeat(allowed, blocked)
```

Each call to `Warrant`, `Rebuttal`, or `Defeater` returns a `*Rule` reference. Pass these references to `Defeat(from, to)` to declare which rule defeats which.

## ContextBuilderFunc

You provide the bridge between `gin.Context` and `RouteContext`:

```go
func buildCtx(c *gin.Context) *route.RouteContext {
    return &route.RouteContext{
        User:     getUserFromJWT(c),  // your JWT logic
        ClientIP: c.ClientIP(),
        Method:   c.Request.Method,
        Path:     c.Request.URL.Path,
        Headers:  extractHeaders(c),
    }
}
```
