# pkg/policy

**Stop nesting if-else for access control. Declare rules and relationships.**

Policy judgment built on toulmin defeats graph. Authentication, authorization, IP blocking, rate limiting — all as declarative rule relationships. Framework independent (net/http).

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

mux := http.NewServeMux()
mux.Handle("/admin/users", policy.Guard(g, buildCtx)(adminHandler))
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
| Signature | `func(http.Handler) http.Handler` | `func(http.Handler) http.Handler` |

Both deny when `verdict <= 0`.

## ContextFunc

```go
func buildCtx(r *http.Request) *policy.RequestContext {
    return &policy.RequestContext{
        User:     getUserFromJWT(r),  // your domain User type
        ClientIP: r.RemoteAddr,
        Headers:  extractHeaders(r),
    }
}
```

## Web Framework Integration

```go
// net/http
mux.Handle("/admin", policy.Guard(g, buildCtx)(handler))

// Gin
r.GET("/admin", func(c *gin.Context) {
    rc := buildCtxFromGin(c)
    results, _ := g.Evaluate(nil, rc)
    if results[0].Verdict <= 0 {
        c.AbortWithStatusJSON(403, gin.H{"error": "forbidden"})
        return
    }
    c.Next()
})

// Chi
r.Use(policy.Guard(g, buildCtx))

// Echo
e.Use(echo.WrapMiddleware(policy.Guard(g, buildCtx)))
```
