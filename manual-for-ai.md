# toulmin — Manual for AI Agents

Toulmin argumentation-based rule engine for Go. Rules are Go functions. Engine builds defeats graph and computes verdict via h-Categoriser.

## Framework Packages

| Package | Description | README |
|---|---|---|
| `pkg/toulmin` | Core engine (Graph, Rule, Evaluate, h-Categoriser) | `pkg/toulmin/README.md` |
| `pkg/policy` | Policy judgment (auth, IP, rate limit, net/http Guard) | `pkg/policy/README.md` |
| `pkg/state` | State transition (Machine, Mermaid diagram) | `pkg/state/README.md` |
| `pkg/approve` | Approval workflow (multi-step Flow) | `pkg/approve/README.md` |
| `pkg/price` | Price judgment (DiscountBacking, Pricer) | `pkg/price/README.md` |
| `pkg/feature` | Feature flags (Flags, rollout, net/http Require/Inject) | `pkg/feature/README.md` |
| `pkg/moderate` | Content moderation (Classifier, 3-level action) | `pkg/moderate/README.md` |

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

Spec structs must implement `SpecName()` (returns identifier for ruleID) and `Validate()` (validates fields at registration time). `nil` specs is allowed for rules that don't need criteria. Func fields in Spec structs are forbidden — `Validate()` rejects them.

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
results, err = g.Evaluate(ctx, toulmin.EvalOption{Method: toulmin.Recursive})             // recursive h-Categoriser
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
| `Matrix` (default) | Matrix multiplication verdict computation |
| `Recursive` | Proven recursive h-Categoriser traversal |

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
    Specs     any           `json:"specs,omitempty"`
    Duration  time.Duration `json:"duration,omitempty"`
}
```

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
        {Name: "admin allowed",  Ctx: &Ctx{Role: "admin"},  Expect: toulmin.VerdictAbove(0)},
        {Name: "blocked IP",     Ctx: &Ctx{IP: "blocked"},  Expect: toulmin.VerdictAtMost(0)},
        {Name: "unauthenticated", Ctx: &Ctx{User: nil},     Expect: toulmin.NoResult},
    })
}
```

### TestCase

```go
type TestCase struct {
    Name   string      // sub-test name
    Ctx    Context     // passed to Evaluate (use NewContext())
    Option EvalOption  // zero value for defaults
    Expect Expectation // verdict assertion
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
| Func field in Spec struct | `Validate()` rejects func fields — use plain data fields only |

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
