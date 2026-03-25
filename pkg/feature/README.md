# pkg/feature

**Stop nesting if-else for feature flags. Declare conditions, exceptions, and rollouts.**

Feature flag framework built on toulmin defeats graph. Toggle, percentage rollout, regional targeting, exception handling — all declarative. No SaaS dependency. Framework independent (net/http).

## Install

```go
import "github.com/park-jun-woo/toulmin/pkg/feature"
```

## Quick Start

```go
g := toulmin.NewGraph("feature:dark-mode")
beta := g.Rule(feature.IsBetaUser)
legacy := g.Counter(feature.IsLegacyBrowser)
internal := g.Except(feature.IsInternalStaff)
legacy.Attacks(beta)
internal.Attacks(legacy)

flags := feature.NewFlags()
flags.Register("dark-mode", g)

enabled, _ := flags.IsEnabled("dark-mode", ctx)
result, _ := flags.Evaluate("dark-mode", ctx, toulmin.EvalOption{Trace: true})
active, _ := flags.List(ctx)
```

## Rules

| Rule | Backing | Description |
|---|---|---|
| `IsBetaUser` | nil | Attributes["beta"] is true |
| `IsInternalStaff` | nil or func(any)string | Internal staff check |
| `IsRegion` | string | User region matches backing |
| `HasAttribute` | [2]any | Attributes[key] == value |
| `IsLegacyBrowser` | nil | Attributes["legacy_browser"] is true |
| `IsUserInPercentage` | float64 | Deterministic hash rollout (0.3 = 30%) |

## Flags

```go
flags.IsEnabled(name, ctx)                                    // bool
flags.Evaluate(name, ctx, toulmin.EvalOption{Trace: true})    // FeatureResult with trace
flags.List(ctx)                                               // all enabled feature names
```

## Middleware (net/http)

```go
// Require — returns 404 if feature is disabled
mux.Handle("/v2/checkout", feature.Require(flags, "new-checkout", buildCtx)(handler))

// Inject — stores enabled features in request context
mux.Handle("/", feature.Inject(flags, buildCtx)(handler))

// Retrieve injected features
active := feature.ActiveFeatures(r)
```

## Web Framework Integration

```go
// Gin
r.GET("/v2/checkout", func(c *gin.Context) {
    ctx := buildCtxFromGin(c)
    enabled, _ := flags.IsEnabled("new-checkout", ctx)
    if !enabled {
        c.AbortWithStatus(404)
        return
    }
    c.Next()
})

// Chi
r.Use(feature.Require(flags, "new-checkout", buildCtx))

// Echo
e.Use(echo.WrapMiddleware(feature.Require(flags, "new-checkout", buildCtx)))
```

## Deterministic Rollout

`IsUserInPercentage` uses FNV-1a hash — same user always gets the same result. No randomness.
