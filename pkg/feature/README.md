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

enabled, _ := flags.IsEnabled("dark-mode", uctx)
result, _ := flags.EvaluateTrace("dark-mode", uctx)
active, _ := flags.List(uctx)
```

## Rules

| Rule | Spec | Description |
|---|---|---|
| `IsBetaUser` | nil | Attributes["beta"] is true |
| `IsInternalStaff` | nil | Attributes["internal"] is true |
| `IsRegion` | *RegionSpec | User region matches spec.Region |
| `HasAttribute` | *AttributeSpec | Attributes[spec.Key] == spec.Value |
| `IsLegacyBrowser` | nil | Attributes["legacy_browser"] is true |
| `IsUserInPercentage` | *PercentageSpec | Deterministic hash rollout (spec.Percentage, e.g. 0.3 = 30%) |

### Spec Types

```go
type RegionSpec struct {
    Region string // target region code ("KR", "US", etc.)
}

type AttributeSpec struct {
    Key   string // attribute key
    Value any    // attribute value to match
}

type PercentageSpec struct {
    Percentage float64 // rollout percentage (0.0 ~ 1.0)
}
```

## Flags

```go
flags.IsEnabled(name, uctx)           // bool
flags.EvaluateTrace(name, uctx)       // FeatureResult with trace
flags.List(uctx)                       // all enabled feature names
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
    uctx := buildCtxFromGin(c)
    enabled, _ := flags.IsEnabled("new-checkout", uctx)
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
