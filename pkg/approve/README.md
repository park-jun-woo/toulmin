# pkg/approve

**Stop nesting if-else for approval logic. Declare steps, rules, and exceptions.**

Multi-step approval workflow built on toulmin defeats graph. Each step is a graph, exceptions are defeat edges, audit trail is built-in.

## Install

```go
import "github.com/park-jun-woo/toulmin/pkg/approve"
```

## Quick Start

```go
g := toulmin.NewGraph("expense:manager")
budget := g.Warrant(approve.IsUnderBudget, nil, 1.0)
frozen := g.Rebuttal(approve.IsBudgetFrozen, nil, 1.0)
urgent := g.Defeater(approve.IsUrgent, nil, 1.0)
g.Defeat(frozen, budget)
g.Defeat(urgent, frozen)

f := approve.NewFlow("expense").AddStep("manager", g)

result, _ := f.Evaluate(req, func(step string) *approve.ApprovalContext {
    return &approve.ApprovalContext{
        Approver: manager,
        Budget:   deptBudget,
        OrgTree:  orgTree,
    }
})
// result.Approved: true/false
// result.Steps[0].Trace: audit trail with backing values
```

## Rules

| Rule | Backing | Description |
|---|---|---|
| `IsDirectManager` | nil | Approver is requester's direct manager |
| `IsUnderBudget` | nil | Amount <= remaining budget |
| `IsBudgetFrozen` | nil | Budget is frozen (Rebuttal) |
| `HasApprovalRole` | string | Approver role matches backing |
| `IsAboveLevel` | int | Approver level >= backing |
| `IsSmallAmount` | float64 | Amount <= backing threshold |
| `IsUrgent` | nil | Request metadata has urgent=true |
| `IsCEOOverride` | nil | Approver role is "ceo" (Defeater) |

## Flow

```go
f := approve.NewFlow("expense").
    AddStep("manager", mgrGraph).
    AddStep("finance", finGraph)
```

- Steps run sequentially
- All steps must pass (verdict > 0)
- Stops at first rejection
- Each step gets its own `ApprovalContext` via `StepContextFunc`

## FlowResult

```go
type FlowResult struct {
    Approved bool
    Steps    []StepResult  // verdict + trace per step
}
```
