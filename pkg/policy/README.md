# pkg/policy

**Stop nesting if-else for access control. Declare rules and relationships.**

Policy judgment built on toulmin defeats graph. Authentication, authorization, IP blocking, rate limiting — all as declarative rule relationships. Works with Gin.

User is `any` — the framework does not impose a concrete User type. Field access is done via backing (extraction functions).

## Install

```go
import "github.com/park-jun-woo/toulmin/pkg/policy"
```

## Quick Start

```go
roleFunc := func(u any) string { return u.(*MyUser).Role }

blocklist := &policy.IPListBacking{Purpose: "blocklist", Check: isBlocked}
whitelist := &policy.IPListBacking{Purpose: "whitelist", Check: isWhitelisted}

g := toulmin.NewGraph("admin:users")
auth := g.Warrant(policy.IsAuthenticated, nil, 1.0)
admin := g.Warrant(policy.IsInRole, &policy.RoleBacking{Role: "admin", RoleFunc: roleFunc}, 1.0)
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
| `IsInRole` | *RoleBacking | User role matches backing.Role (via RoleFunc) |
| `IsOwner` | *OwnerBacking | User ID matches resource owner (via extraction funcs) |
| `IsIPInList` | *IPListBacking | Client IP in list (purpose + check) |
| `IsRateLimited` | nil | Client IP is rate limited |
| `HasHeader` | string | Named header exists and is non-empty |

### Backing Types

```go
type RoleBacking struct {
    Role     string
    RoleFunc func(any) string  // extracts role from domain User
}

type OwnerBacking struct {
    UserIDFunc     func(any) string  // extracts user ID
    ResourceIDFunc func(any) string  // extracts resource owner ID
}

type IPListBacking struct {
    Purpose string             // "blocklist", "whitelist"
    Check   func(string) bool
}
```

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
        User:     getUserFromJWT(c),  // your domain User type
        ClientIP: c.ClientIP(),
        Headers:  extractHeaders(c),
    }
}
```
