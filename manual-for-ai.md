# toulmin — Manual for AI Agents

Toulmin argumentation-based rule engine for Go. Rules are Go functions. Engine builds defeats graph and computes verdict via h-Categoriser.

## Framework Packages

| Package | Description | README |
|---|---|---|
| `pkg/toulmin` | Core engine (Graph, Rule, Evaluate, h-Categoriser) | `pkg/toulmin/README.md` |
| `pkg/policy` | Policy judgment (auth, IP, rate limit, net/http Guard) | `pkg/policy/README.md` |
| `pkg/state` | State transition (Machine, Mermaid diagram) | `pkg/state/README.md` |
| `pkg/approve` | Approval workflow (multi-step Flow) | `pkg/approve/README.md` |
| `pkg/price` | Price judgment (DiscountSpec, Pricer) | `pkg/price/README.md` |
| `pkg/feature` | Feature flags (Flags, rollout, net/http Require/Inject) | `pkg/feature/README.md` |
| `pkg/moderate` | Content moderation (Classifier, 3-level action) | `pkg/moderate/README.md` |
| `pkg/tangl` | Markdown-based policy language (parser, validate) | — |

## How to Navigate

1. Read `codebook.yaml` — project vocabulary
2. Full read only the files you need, then work

---

## Core Concepts

### Rule

```go
func(ctx Context, specs Specs) (bool, any)
```

Returns `(judgment, evidence)`. `ctx` is a Context interface with `Get(key string) (any, bool)` and `Set(key string, value any)`. `specs` receives judgment criteria from graph declaration via `.With()`.

### Spec Interface

```go
type Spec interface {
    SpecName() string
    Validate() error
}
```

Spec structs must implement `SpecName()` (returns identifier for ruleID) and `Validate()` (validates fields at registration time). `nil` specs is allowed for rules that don't need criteria. Func fields in Spec structs are forbidden — `With()` rejects them via internal validation.

### Ground vs Spec

| | Ground | Spec |
|---|---|---|
| What | Facts about judgment target | Judgment criteria |
| When | Per request (runtime) | Fixed at declaration |
| Passed by | `ctx.Set(key, value)` before `Evaluate(ctx)` | `g.Rule(fn).With(spec)` |
| Example | User, IP, request context | Threshold, role name, config |

### Strength

| Strength | Effect |
|---|---|
| Strict | Rejects all incoming attack edges |
| Defeasible | Accepts incoming attack edges |
| Defeater | Outgoing attack edges only, no own verdict |

### Defeats Graph

- **Rule**: node that can be attacked
- **Counter**: node that attacks a rule (has own conclusion)
- **Except**: node that attacks without own conclusion

Distinction is in **graph position** (defeat edges), not function signature.

### Verdict

```
raw(a) = w(a) / (1 + Σ raw(attackers))     [0, 1]
verdict(a) = 2 × raw(a) - 1                [-1, 1]
```

`+1.0` confirmed, `0.0` undecided, `-1.0` fully rebutted. `> 0` warrant prevails, `< 0` rebuttal prevails.

### Qualifier

`0.0–1.0` float. Initial weight per rule. Default `1.0`.

### Rule Identity

`ruleID = funcID + "#" + spec` (non-nil spec). `ruleID = funcID` (nil spec). Same function + different spec = different rule.

---

## Writing Rules

```go
func CheckOneFileOneFunc(ctx toulmin.Context, specs toulmin.Specs) (bool, any) {
    gf, _ := ctx.Get("file")
    f := gf.(*FileGround)
    if len(f.Funcs) > 1 {
        return true, &Evidence{Got: len(f.Funcs), Expected: 1}
    }
    return false, nil
}
```

---

## Graph API

Rule/Counter/Except return `*Rule` reference with builder pattern. Attacks is a method on `*Rule`.

```go
g := toulmin.NewGraph("voting")
w := g.Rule(IsAdult)                           // spec/qualifier optional via builder
r := g.Counter(HasCriminalRecord)
r.Attacks(w)

ctx := toulmin.NewContext()
results, err := g.Evaluate(ctx)                                                          // default (matrix)
results, err = g.Evaluate(ctx, toulmin.EvalOption{Method: toulmin.Recursive})             // not yet implemented — returns error
results, err = g.Evaluate(ctx, toulmin.EvalOption{Trace: true})                           // with trace
results, err = g.Evaluate(ctx, toulmin.EvalOption{Duration: true})                        // with duration (trace auto-enabled)
```

### EvalOption

```go
type EvalOption struct {
    Method   EvalMethod // Matrix (default) | Recursive
    Trace    bool       // collect TraceEntry per warrant
    Duration bool       // measure per-rule execution time (auto-enables Trace)
}
```

| Method | Description |
|---|---|
| `Matrix` (default) | lazy recursive h-Categoriser (default) |
| `Recursive` | Not yet implemented — returns error |

### EvalResult / TraceEntry

```go
type EvalResult struct {
    Name     string       `json:"name"`
    Verdict  float64      `json:"verdict"`
    Evidence any          `json:"evidence,omitempty"`
    Trace    []TraceEntry `json:"trace"`
}

type TraceEntry struct {
    Name      string        `json:"name"`
    Role      string        `json:"role"`
    Activated bool          `json:"activated"`
    Qualifier float64       `json:"qualifier"`
    Evidence  any           `json:"evidence,omitempty"`
    Specs     Specs         `json:"specs,omitempty"`
    Duration  time.Duration `json:"duration,omitempty"`
}
```

### Node Events (Run)

`Run` is a pre-judge → post-act variant of `Evaluate`. It evaluates the whole graph
(full pass, Trace/Duration forced off), fires each node's event handler, and Runs any
sub-graph declared on an Active node (execution composition, below). Every node fires
exactly one event:

| Event | Condition | Meaning |
|---|---|---|
| `Inactive` | func `false` | rule did not apply |
| `Active` | func `true` && `verdict > 0` | applied and prevailed |
| `Defeated` | func `true` && `verdict <= 0` | applied but defeated (verdict == 0 counts as Defeated) |

```go
g.Rule(Authenticate).
    OnActive(func(ctx toulmin.Context, ev toulmin.NodeEvent, view toulmin.RunView) error { return audit(ev) }).
    OnDefeated(func(ctx toulmin.Context, ev toulmin.NodeEvent, view toulmin.RunView) error { return deny(ev) })

results, view, err := g.Run(ctx) // []EvalResult, RunView, error

type NodeEvent struct {
    Name     string        // node short name
    Role     string        // "rule" | "counter" | "except"
    Type     NodeEventType // Inactive | Active | Defeated
    Verdict  float64       // 0.0 when Inactive
    Evidence any           // evidence from the rule function
}

type RunView interface {           // read-only snapshot of every node's final event
    All() []NodeEvent              // all nodes, registration order (Inactive included)
    Get(name string) (NodeEvent, bool) // lookup by short name
    Attackers(name string) []NodeEvent // final events of nodes that attacked name
}
```

Handlers fire in rule registration order (deterministic). Before any handler fires, `Run`
builds one immutable `RunView` snapshot of every node's final event and shares it with all
handlers, so a handler can read the whole graph's final state (audit, explanation, gradient
thresholds via `view.Get(...).Verdict`) — mutating `ctx` never changes the `view`. A handler
error or panic stops `Run` immediately and is returned together with that pre-dispatch
`RunView`. Nodes without a matching handler pass through silently. `Evaluate` is unchanged
and fires no handlers (stays idempotent).

### Execution Composition (graph-of-graphs)

`(r *Rule) Run(g *Graph) *Rule` declares an **execution edge**: when this node is `Active`,
the sub-graph `g` is Run with the *same* `ctx`. It is the execution counterpart of `Attacks`
(a defeat edge). `Attacks` composes **judgment** — an attacker's verdict flows *up* into its
target; `Run` composes **execution** — an Active node drives a child graph *down*, with its
verdict isolated.

```go
notify := buildNotifyGraph()
g.Rule(OrderPlaced).
    OnActive(func(ctx toulmin.Context, ev toulmin.NodeEvent, v toulmin.RunView) error { return log(ev) }).
    Run(notify)   // when OrderPlaced is Active, Run notify
```

- **Active-only.** A sub-graph Runs only for `Active` nodes; `Inactive`/`Defeated` never trigger one.
- **Handler first, then sub-graph.** For an Active node, its `OnActive` handler fires before its sub-graph is Run. `OnActive` and `Run(g)` coexist.
- **ctx flows down, verdict isolated.** The same mutable `ctx` is shared with the sub-graph (side effects propagate); the sub-graph's verdicts are *not* merged into the parent — only errors propagate, wrapped as `run "node" → "subgraph": ...`.
- **Each level gets its own view.** The `RunView` passed to a sub-graph's handlers snapshots that sub-graph, not the parent.
- **DAG, enforced.** Execution composition must be acyclic. `Run` rejects a static cycle once at the top-level entry via `detectRunCycle` — a 3-color DFS over `RunGraph` edges keyed by `*Graph` identity; a shared sub-graph reached by two paths (diamond) is legal. A runtime depth guard (`runMaxDepth = 64`) backstops runaway composition (`run depth exceeded 64`).
- `Run(nil)` is a registration error and panics.

The Run handler + `RunView` family above is available in the Go, TypeScript, and Python ports.

### Same Function, Different Spec

```go
g := toulmin.NewGraph("limits")
w1 := g.Rule(CheckThreshold).With(&ThresholdSpec{Max: 100})
w2 := g.Rule(CheckThreshold).With(&ThresholdSpec{Max: 200}).Qualifier(0.8)
r := g.Counter(HasExemption).With(&ExemptionSpec{Type: "vip"})
r.Attacks(w1)
```

---

## Testing Helper

`RunCases` eliminates boilerplate for table-driven policy tests.

```go
func TestAccessPolicy(t *testing.T) {
    g := buildAccessGraph()
    toulmin.RunCases(t, g, []toulmin.TestCase{
        {Name: "admin allowed",  Context: adminCtx,  Expect: toulmin.VerdictAbove(0)},
        {Name: "blocked IP",     Context: blockedCtx,  Expect: toulmin.VerdictAtMost(0)},
        {Name: "unauthenticated", Context: unauthCtx,     Expect: toulmin.NoResult},
    })
}
```

### TestCase

```go
type TestCase struct {
    Name    string      // sub-test name
    Context Context     // passed to Evaluate (use NewContext())
    Option  EvalOption  // zero value for defaults
    Expect  Expectation // verdict assertion
}
```

### Expectation

| Function | Condition |
|---|---|
| `VerdictAbove(v)` | verdict > v |
| `VerdictAtMost(v)` | verdict <= v |
| `VerdictBetween(lo, hi)` | lo < verdict <= hi |
| `NoResult` | no active warrants (empty results) |

---

## Cycle Detection

Detected at evaluation time via DFS.

---

## Commands

```bash
toulmin evaluate                              # run example
```

---

## Evaluation Flow

```
0. Cycle detection (DFS) → error if cycle found
1. Each rule node → run func(ctx, specs) → false? skip
2. If true → traverse attackers recursively
3. Each attacker: func → false? contributes 0 → true? recurse deeper
4. h-Categoriser: raw(a) = w(a) / (1 + Σ raw(attackers)), verdict = 2*raw - 1
5. Cache per ruleID — each rule runs at most once
6. Only reachable rules executed. Specs passed as 2nd arg.
```

---

## Common Mistakes

| Mistake | Fix |
|---|---|
| Rule func wrong signature | `func(ctx Context, specs Specs) (bool, any)` |
| Chaining calls | Rule/Counter/Except return `*Rule` with builder pattern — spec/qualifier are optional |
| Attacks without registration | Must Rule/Counter/Except first to get `*Rule` |
| Verdict 0.0 as allow/deny | 0.0 = undecided — threshold is framework's decision |
| Confusing context and spec | ctx = per-request facts via Get/Set, spec = fixed criteria at declaration |
| Forgetting spec | Use `nil` when no spec needed |
| Func field in Spec struct | `With()` rejects func fields via internal validation — use plain data fields only |

### Spec Replaces Closures

Same function + different spec values — no closure factories needed. `ruleID` distinguishes them.

```go
g := toulmin.NewGraph("example")
r1 := g.Counter(HasRole).With(&RoleSpec{Role: "admin"})   // ruleID = "HasRole#admin"
r2 := g.Counter(HasRole).With(&RoleSpec{Role: "editor"})  // ruleID = "HasRole#editor"
r1.Attacks(someRule)
```

### Verdict 0.0 Threshold

- **Security**: `verdict <= 0` → deny
- **Moderation**: `verdict < 0.3` → flag
- **Feature flags**: `verdict > 0` → enabled

Engine computes verdict. Framework interprets it.
