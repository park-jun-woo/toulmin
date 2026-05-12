# pkg/moderate

**Stop nesting if-else for content moderation. Declare rules, exceptions, and context.**

Content moderation framework built on toulmin defeats graph. Hate speech, spam, NSFW detection as rebuttals. News context, education, trusted users as defeaters. AI classifiers plug in via spec. Audit trail is built-in. Framework independent (net/http).

## Install

```go
import "github.com/park-jun-woo/toulmin/pkg/moderate"
```

## Quick Start

```go
classifier := myClassifier{}

g := toulmin.NewGraph("post:publish")
verified := g.Rule(moderate.IsVerifiedUser)
hate := g.Counter(moderate.ContainsHateSpeech).With(&moderate.ClassifierSpec{Classifier: classifier})
news := g.Except(moderate.IsNewsContext)
hate.Attacks(verified)
news.Attacks(hate)

mod := moderate.NewModerator(g)
result, _ := mod.Review(content, ctx)
// result.Action: "allow", "flag", or "block"
// result.Trace: full audit trail
```

## 3-Level Action

| Action | Verdict | HTTP |
|---|---|---|
| `allow` | > 0.3 | 200 (next) |
| `flag` | 0 < v <= 0.3 | 202 (manual review) |
| `block` | <= 0 | 403 |

## Rules

| Rule | Spec | Description |
|---|---|---|
| `IsVerifiedUser` | nil | Author is verified |
| `IsTrustedUser` | *TrustScoreSpec | Author trust score >= spec.MinScore |
| `ContainsHateSpeech` | *ClassifierSpec | AI score > 0.8 |
| `ContainsSpam` | *ClassifierSpec | AI score > 0.7 |
| `ContainsNSFW` | *ClassifierSpec | AI score > 0.8 |
| `IsNewsContext` | nil | Channel type is "news" |
| `IsEducational` | nil | Channel type is "education" |
| `IsAdultChannel` | nil | Channel is age-gated |
| `HasMinPosts` | *MinPostsSpec | Author post count >= spec.MinPosts |

### Spec Types

```go
type ClassifierSpec struct {
    Classifier Classifier // AI classifier implementation
}

type TrustScoreSpec struct {
    MinScore float64 // minimum trust score threshold
}

type MinPostsSpec struct {
    MinPosts int // minimum number of posts required
}
```

## Classifier Interface

```go
type Classifier interface {
    Predict(text string, category string) float64
}
```

Plug in any AI model, LLM, keyword matcher, or external API.

## Middleware (net/http)

```go
mux.Handle("/posts", moderate.Guard(mod, extractPost)(handler))
// Block → 403, Flag → 202, Allow → next
```

## Web Framework Integration

```go
// Gin
r.POST("/posts", func(c *gin.Context) {
    content, ctx := extractPostFromGin(c)
    result, _ := mod.Review(content, ctx)
    switch result.Action {
    case moderate.ActionBlock:
        c.AbortWithStatusJSON(403, gin.H{"error": "blocked"})
    case moderate.ActionFlag:
        c.AbortWithStatusJSON(202, gin.H{"status": "flagged"})
    default:
        c.Next()
    }
})

// Chi
r.Use(moderate.Guard(mod, extractPost))

// Echo
e.Use(echo.WrapMiddleware(moderate.Guard(mod, extractPost)))
```
