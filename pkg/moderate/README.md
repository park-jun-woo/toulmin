# pkg/moderate

**Stop nesting if-else for content moderation. Declare rules, exceptions, and context.**

Content moderation framework built on toulmin defeats graph. Hate speech, spam, NSFW detection as rebuttals. News context, education, trusted users as defeaters. AI classifiers plug in via backing. Audit trail is built-in.

## Install

```go
import "github.com/park-jun-woo/toulmin/pkg/moderate"
```

## Quick Start

```go
classifier := myClassifier{}

g := toulmin.NewGraph("post:publish")
verified := g.Warrant(moderate.IsVerifiedUser, nil, 1.0)
hate := g.Rebuttal(moderate.ContainsHateSpeech, classifier, 1.0)
news := g.Defeater(moderate.IsNewsContext, nil, 1.0)
g.Defeat(hate, verified)
g.Defeat(news, hate)

mod := moderate.NewModerator(g)
result, _ := mod.Review(content, ctx)
// result.Action: "allow", "flag", or "block"
// result.Trace: full audit trail
```

## 3-Level Action

| Action | Verdict | HTTP |
|---|---|---|
| `allow` | > 0.3 | c.Next() |
| `flag` | 0 < v <= 0.3 | 202 (manual review) |
| `block` | <= 0 | 403 |

## Rules

| Rule | Backing | Description |
|---|---|---|
| `IsVerifiedUser` | nil | Author is verified |
| `IsTrustedUser` | float64 | Author trust score >= backing |
| `ContainsHateSpeech` | Classifier | AI score > 0.8 |
| `ContainsSpam` | Classifier | AI score > 0.7 |
| `ContainsNSFW` | Classifier | AI score > 0.8 |
| `IsNewsContext` | nil | Channel type is "news" |
| `IsEducational` | nil | Channel type is "education" |
| `IsAdultChannel` | nil | Channel is age-gated |
| `HasMinPosts` | int | Author post count >= backing |

## Classifier Interface

```go
type Classifier interface {
    Predict(text string, category string) float64
}
```

Plug in any AI model, LLM, keyword matcher, or external API.

## Gin Middleware

```go
r.POST("/posts", moderate.Guard(mod, extractPost), handler)
// Block → 403, Flag → 202, Allow → next
```
