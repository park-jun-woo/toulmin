# pkg/feature

**Stop nesting if-else for feature flags. Declare conditions, exceptions, and rollouts.**

Feature flag framework built on toulmin defeats graph. Toggle, percentage rollout, regional targeting, exception handling — all declarative. No SaaS dependency.

## Install

```go
import "github.com/park-jun-woo/toulmin/pkg/feature"
```

## Quick Start

```go
g := toulmin.NewGraph("feature:dark-mode")
beta := g.Warrant(feature.IsBetaUser, nil, 1.0)
legacy := g.Rebuttal(feature.IsLegacyBrowser, nil, 1.0)
internal := g.Defeater(feature.IsInternalStaff, nil, 1.0)
g.Defeat(legacy, beta)
g.Defeat(internal, legacy)

flags := feature.NewFlags()
flags.Register("dark-mode", g)

enabled, _ := flags.IsEnabled("dark-mode", ctx)
result, _ := flags.EvaluateTrace("dark-mode", ctx)
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
flags.IsEnabled(name, ctx)       // bool
flags.EvaluateTrace(name, ctx)   // FeatureResult with trace
flags.List(ctx)                  // all enabled feature names
```

## Gin Middleware

```go
r.GET("/v2/checkout", feature.Require(flags, "new-checkout", buildCtx), handler)
r.Use(feature.Inject(flags, buildCtx))  // sets "features" in gin.Context
```

## Deterministic Rollout

`IsUserInPercentage` uses FNV-1a hash — same user always gets the same result. No randomness.
