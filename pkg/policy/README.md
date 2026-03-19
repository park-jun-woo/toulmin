# pkg/policy

**Stop nesting if-else for access control. Declare rules and relationships.**

Policy judgment built on toulmin defeats graph. Authentication, authorization, IP blocking, rate limiting — all as declarative rule relationships. Works with Gin.

## Install

```go
import "github.com/park-jun-woo/toulmin/pkg/policy"
```

## Quick Start

```go
blocklist := &policy.IPListBacking{Purpose: "blocklist", Check: isBlocked}
whitelist := &policy.IPListBacking{Purpose: "whitelist", Check: isWhitelisted}

g := toulmin.NewGraph("admin:users")
auth := g.Warrant(policy.IsAuthenticated, nil, 1.0)
admin := g.Warrant(policy.IsInRole, "admin", 1.0)
blocked := g.Rebuttal(policy.IsIPInList, blocklist, 1.0)
allowed := g.Defeater(policy.IsIPInList, whitelist, 1.0)
g.Defeat(blocked, auth)
g.Defeat(allowed, blocked)

r := gin.Default()
r.GET("/admin/users", policy.Guard(g, buildCtx), handler)
```

## Rules

`func(claim, ground, backing) (bool, any)` — ground is per-request facts, backing is fixed criteria.

| Rule | Backing | Description |
|---|---|---|
| `IsAuthenticated` | nil | User is not nil |
| `IsInRole` | string | User role matches backing |
| `IsOwner` | nil | User ID matches resource owner |
| `IsIPInList` | *IPListBacking | Client IP in list (purpose + check) |
| `IsRateLimited` | nil | Client IP is rate limited |
| `HasHeader` | string | Named header exists and is non-empty |

### IPListBacking

```go
type IPListBacking struct {
    Purpose string            // "blocklist", "whitelist"
    Check   func(string) bool
}
```

Same function, different backing — purpose distinguishes intent, check performs lookup.

## Guard vs GuardDebug

| | Guard | GuardDebug |
|---|---|---|
| Evaluation | `Evaluate` | `EvaluateTrace` |
| Headers | None | `X-Policy-Verdict`, `X-Policy-Trace` |
| 403 body | `{"error":"forbidden"}` | `{"error":"forbidden","trace":"..."}` |

Both deny when `verdict <= 0`.

## ContextBuilderFunc

```go
func buildCtx(c *gin.Context) *policy.RequestContext {
    return &policy.RequestContext{
        User:     getUserFromJWT(c),
        ClientIP: c.ClientIP(),
        Headers:  extractHeaders(c),
    }
}
```
