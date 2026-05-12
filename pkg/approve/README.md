# pkg/approve

**Stop nesting if-else for approval logic. Declare steps, rules, and exceptions.**

Multi-step approval workflow built on toulmin defeats graph. Each step is a graph, exceptions are defeat edges, audit trail is built-in.

Approver is `any` — the framework does not impose a concrete Approver type. Field access is done via `ApprovalContext` fields set in `ctx`.

## Install

```go
import "github.com/park-jun-woo/toulmin/pkg/approve"
```

## Quick Start

```go
g := toulmin.NewGraph("expense:manager")
budget := g.Rule(approve.IsUnderBudget)
role := g.Rule(approve.HasApprovalRole).With(&approve.ApproverSpec{Role: "manager"})
frozen := g.Counter(approve.IsBudgetFrozen)
urgent := g.Except(approve.IsUrgent)
frozen.Attacks(budget)
urgent.Attacks(frozen)

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

| Rule | Spec | Description |
|---|---|---|
| `IsDirectManager` | nil | Approver is requester's direct manager (via ctx "orgTree", "approverID", "requesterID") |
| `IsUnderBudget` | nil | Amount <= remaining budget |
| `IsBudgetFrozen` | nil | Budget is frozen (Rebuttal) |
| `HasApprovalRole` | *ApproverSpec | Approver role matches spec.Role (via ctx "approverRole") |
| `IsAboveLevel` | *ApproverSpec | Approver level >= spec.Level (via ctx "approverLevel") |
| `IsSmallAmount` | *ThresholdSpec | Amount <= spec.Max |
| `IsUrgent` | nil | Request metadata has urgent=true |
| `IsCEOOverride` | nil | Approver role is "ceo" (via ctx "approverRole") |

### Spec Types

```go
type ApproverSpec struct {
    Role  string // role to match (for HasApprovalRole)
    Level int    // minimum level (for IsAboveLevel)
}

type ThresholdSpec struct {
    Max float64 // amount threshold
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
