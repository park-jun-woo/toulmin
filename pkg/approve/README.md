# pkg/approve

**Stop nesting if-else for approval logic. Declare steps, rules, and exceptions.**

Multi-step approval workflow built on toulmin defeats graph. Each step is a graph, exceptions are defeat edges, audit trail is built-in.

Approver is `any` — the framework does not impose a concrete Approver type. Field access is done via backing (ApproverBacking extraction functions).

## Install

```go
import "github.com/park-jun-woo/toulmin/pkg/approve"
```

## Quick Start

```go
ab := &approve.ApproverBacking{
    IDFunc:    func(a any) string { return a.(*MyApprover).ID },
    RoleFunc:  func(a any) string { return a.(*MyApprover).Role },
    LevelFunc: func(a any) int { return a.(*MyApprover).Level },
}

g := toulmin.NewGraph("expense:manager")
budget := g.Warrant(approve.IsUnderBudget, nil, 1.0)
frozen := g.Rebuttal(approve.IsBudgetFrozen, nil, 1.0)
urgent := g.Defeater(approve.IsUrgent, nil, 1.0)
g.Defeat(frozen, budget)
g.Defeat(urgent, frozen)

f := approve.NewFlow("expense").AddStep("manager", g)

result, _ := f.Evaluate(req, func(step string) *approve.ApprovalContext {
    return &approve.ApprovalContext{
        Approver: myApprover,  // your domain Approver type
        Budget:   deptBudget,
        OrgTree:  orgTree,
    }
})
```

## Rules

| Rule | Backing | Description |
|---|---|---|
| `IsDirectManager` | *ApproverBacking | Approver is requester's direct manager (via IDFunc) |
| `IsUnderBudget` | nil | Amount <= remaining budget |
| `IsBudgetFrozen` | nil | Budget is frozen (Rebuttal) |
| `HasApprovalRole` | *ApproverBacking | Approver role matches backing.Role (via RoleFunc) |
| `IsAboveLevel` | *ApproverBacking | Approver level >= backing.Level (via LevelFunc) |
| `IsSmallAmount` | float64 | Amount <= backing threshold |
| `IsUrgent` | nil | Request metadata has urgent=true |
| `IsCEOOverride` | *ApproverBacking | Approver role is "ceo" (via RoleFunc) |

### ApproverBacking

```go
type ApproverBacking struct {
    Role      string
    Level     int
    IDFunc    func(any) string
    RoleFunc  func(any) string
    LevelFunc func(any) int
}
```

## Flow

```go
f := approve.NewFlow("expense").
    AddStep("manager", mgrGraph).
    AddStep("finance", finGraph)
```

- Steps run sequentially — all must pass (verdict > 0)
- Stops at first rejection
- Each step gets its own `ApprovalContext` via `StepContextFunc`
