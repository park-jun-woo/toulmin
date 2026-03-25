# pkg/policy

**Stop nesting if-else for access control. Declare rules and relationships.**

Policy judgment built on toulmin defeats graph. Authentication, authorization, IP blocking, rate limiting — all as declarative rule relationships. Framework independent (net/http).

User is `any` — the framework does not impose a concrete User type. Field access is done via `RequestContext` fields.

## Install

```go
import "github.com/park-jun-woo/toulmin/pkg/policy"
```

## Quick Start

```go
g := toulmin.NewGraph("admin:users")
auth := g.Rule(policy.IsAuthenticated)
admin := g.Rule(policy.IsInRole).Backing(&policy.RoleBacking{Role: "admin"})
blocked := g.Counter(policy.IsIPInList).Backing(&policy.IPListBacking{Purpose: "blocklist"})
allowed := g.Except(policy.IsIPInList).Backing(&policy.IPListBacking{Purpose: "whitelist"})
blocked.Attacks(auth)
allowed.Attacks(blocked)

mux := http.NewServeMux()
mux.Handle("/admin/users", policy.Guard(g, buildCtx)(adminHandler))
```

## Rules

`func(ctx toulmin.Context, backing toulmin.Backing) (bool, any)` — ctx provides per-request facts via Get/Set, backing is fixed criteria.

| Rule | Backing | Description |
|---|---|---|
| `IsAuthenticated` | nil | User is not nil |
| `IsInRole` | *RoleBacking | User role matches backing.Role (via RequestContext.Role) |
| `IsOwner` | *OwnerBacking | User ID matches resource owner (via RequestContext.UserID/ResourceOwner) |
| `IsIPInList` | *IPListBacking | Client IP in list (via RequestContext.IPBlocked) |
| `IsRateLimited` | nil | Client IP is rate limited |
| `HasHeader` | string | Named header exists and is non-empty |

### Backing Types

All backing types implement `BackingName() string` and `Validate() error`. Func fields are forbidden.

```go
type RoleBacking struct {
    Role string
}

type OwnerBacking struct {
    OwnerField string  // field name to match against UserID
}

type IPListBacking struct {
    Purpose string  // "blocklist", "whitelist"
}
```

## Guard vs GuardDebug

| | Guard | GuardDebug |
|---|---|---|
| Evaluation | `Evaluate(ctx)` | `Evaluate(ctx, EvalOption{Trace: true})` |
| Headers | None | `X-Policy-Verdict`, `X-Policy-Trace` |
| 403 body | `{"error":"forbidden"}` | `{"error":"forbidden","trace":"..."}` |
| Signature | `func(http.Handler) http.Handler` | `func(http.Handler) http.Handler` |

Both deny when `verdict <= 0`.

## ContextFunc

```go
func buildCtx(r *http.Request) *policy.RequestContext {
    user := getUserFromJWT(r)
    return &policy.RequestContext{
        User:          user,              // your domain User type
        Role:          user.Role,         // extracted role string
        UserID:        user.ID,           // extracted user ID
        ResourceOwner: getOwnerID(r),     // resource owner ID
        ClientIP:      r.RemoteAddr,
        IPBlocked:     isBlocked(r.RemoteAddr),
        Headers:       extractHeaders(r),
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
    ctx := toulmin.NewContext()
    ctx.Set("rc", rc)
    results, _ := g.Evaluate(ctx)
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
