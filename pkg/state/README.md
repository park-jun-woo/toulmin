# pkg/state

**Stop scattering state guards in if-else. Declare transitions and exceptions.**

State transition framework built on toulmin defeats graph. Each transition is a graph, exceptions are defeat edges, Mermaid diagrams are auto-generated.

## Install

```go
import "github.com/park-jun-woo/toulmin/pkg/state"
```

## Quick Start

```go
g := toulmin.NewGraph("proposal:accept")
current := g.Rule(state.IsCurrentState)
owner := g.Rule(state.IsOwner).Backing(ownerBacking)
expired := g.Counter(state.IsExpired).Backing(expiryFunc)
override := g.Except(isAdminOverride)
expired.Attacks(current)
override.Attacks(expired)

m := state.NewMachine()
m.Add("pending", "accept", "accepted", g)

verdict, _ := m.Can(req, ctx)
// verdict > 0: transition allowed
// verdict <= 0: transition denied
```

## Rules

| Rule | Backing | Description |
|---|---|---|
| `IsCurrentState` | nil | ground.CurrentState == claim.From |
| `IsOwner` | *OwnerBacking | User ID matches resource owner ID |
| `IsExpired` | func(any) time.Time | Resource expiry has passed |

### OwnerBacking

```go
type OwnerBacking struct {
    OwnerIDFunc func(any) string
    UserIDFunc  func(any) string
}
```

## Machine

```go
m := state.NewMachine()
m.Add("pending", "accept", "accepted", acceptGraph)
m.Add("pending", "reject", "rejected", rejectGraph)

verdict, err := m.Can(req, ctx)        // verdict only
result, err := m.CanTrace(req, ctx)    // verdict + trace
```

## Mermaid Diagram

```go
fmt.Println(m.Mermaid())
// stateDiagram-v2
//     pending --> accepted : accept
//     pending --> rejected : reject
```

Generated from the same source as runtime guards — always in sync.
