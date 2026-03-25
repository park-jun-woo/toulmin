# toulmin — Manual for AI Agents

Toulmin argumentation-based rule engine for Go. Rules are Go functions. Engine builds defeats graph and computes verdict via h-Categoriser.

## Framework Packages

| Package | Description | README |
|---|---|---|
| `pkg/toulmin` | Core engine (Graph, Rule, Evaluate, h-Categoriser, ParseYAML, ValidateGraphDef, GenerateGraph) | `pkg/toulmin/README.md` |
| `pkg/analyzer` | Go AST analysis (extract defeat graphs from source) | — |
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
func(ctx Context, backing Backing) (bool, any)
```

Returns `(judgment, evidence)`. `ctx` is a Context interface with `Get(key string) (any, bool)` and `Set(key string, value any)`. `backing` receives judgment criteria from graph declaration.

### Backing Interface

```go
type Backing interface {
    BackingName() string
    Validate() error
}
```

Backing structs must implement `BackingName()` (returns identifier for ruleID) and `Validate()` (validates fields at registration time). `nil` backing is allowed for rules that don't need criteria. Func fields in Backing structs are forbidden — `Validate()` rejects them.

### Ground vs Backing

| | Ground | Backing |
|---|---|---|
| What | Facts about judgment target | Judgment criteria |
| When | Per request (runtime) | Fixed at declaration |
| Passed by | `ctx.Set(key, value)` before `Evaluate(ctx)` | `g.Rule(fn).Backing(b)` |
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

`ruleID = funcID + "#" + backing` (non-nil backing). `ruleID = funcID` (nil backing). Same function + different backing = different rule.

---

## Writing Rules

```go
func CheckOneFileOneFunc(ctx toulmin.Context, backing toulmin.Backing) (bool, any) {
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
w := g.Rule(IsAdult)                           // backing/qualifier optional via builder
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
    Backing   any           `json:"backing,omitempty"`
    Duration  time.Duration `json:"duration,omitempty"`
}
```

### Same Function, Different Backing

```go
g := toulmin.NewGraph("limits")
w1 := g.Rule(CheckThreshold).Backing(&ThresholdBacking{Max: 100})
w2 := g.Rule(CheckThreshold).Backing(&ThresholdBacking{Max: 200}).Qualifier(0.8)
r := g.Counter(HasExemption).Backing(&ExemptionBacking{Type: "vip"})
r.Attacks(w1)
```

### LoadGraph — Dynamic Graph from Definition

Builds a live `*Graph` from a `GraphDef`, function registry, and optional backing registry. No code generation, no recompilation.

```go
funcs := map[string]any{
    "isAuthenticated": isAuthenticated,
    "isIPBlocked":     isIPBlocked,
}
backings := map[string]toulmin.Backing{
    "isIPBlocked": fetchBlocklistFromRedis(),
}

g, err := toulmin.LoadGraph(def, funcs, backings)
ctx := toulmin.NewContext()
results, _ := g.Evaluate(ctx)
```

`GraphDef` can come from YAML, DB, or API — graph structure and backing change without redeployment. Functions stay compiled.

```go
type GraphDef struct {
    Graph   string
    Rules   []GraphRuleDef   // Name, Role, Qualifier
    Defeats []GraphEdgeDef   // From, To
}
```

### YAML Graph Schema

```yaml
graph: <name>              # graph name (required)
rules:                     # rule list (required)
  - name: <rule_name>      # rule name, matches function registry key (required)
    role: <role>           # rule | counter | except (required)
    qualifier: <float>     # 0.0–1.0, default 1.0 if omitted (optional)
defeats:                   # defeat edge list (optional)
  - from: <attacker_name>  # rule name of attacker (must exist in rules)
    to: <target_name>      # rule name of target (must exist in rules)
```

Example:

```yaml
graph: api:access
rules:
  - name: isAuthenticated
    role: rule
  - name: isIPBlocked
    role: counter
  - name: isInternalIP
    role: except
    qualifier: 0.8
defeats:
  - from: isIPBlocked
    to: isAuthenticated
  - from: isInternalIP
    to: isIPBlocked
```

### ParseYAML + ValidateGraphDef

`ParseYAML` parses YAML into `GraphDef` (AST). `ValidateGraphDef` checks graph name, defeat edge references, and cycles.

```go
def, err := toulmin.ParseYAML("policy.yaml")
if err := toulmin.ValidateGraphDef(def); err != nil { /* handle */ }
g, err := toulmin.LoadGraph(def, funcs, backings)
ctx := toulmin.NewContext()
ctx.Set("req", req)
results, _ := g.Evaluate(ctx)
```

### Engine API (legacy)

`Engine.Register(RuleMeta{...})` + `Engine.Evaluate()` — string-based names, still available.

---

## Testing Helper

`RunCases` eliminates boilerplate for table-driven policy tests.

```go
func TestAccessPolicy(t *testing.T) {
    g := buildAccessGraph()
    toulmin.RunCases(t, g, []toulmin.TestCase{
        {Name: "admin allowed",  Ground: &Ctx{Role: "admin"},  Expect: toulmin.VerdictAbove(0)},
        {Name: "blocked IP",     Ground: &Ctx{IP: "blocked"},  Expect: toulmin.VerdictAtMost(0)},
        {Name: "unauthenticated", Ground: &Ctx{User: nil},     Expect: toulmin.NoResult},
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

Detected at evaluation time via DFS. CLI also detects before codegen:

```bash
toulmin graph voting.yaml --check   # validate only
toulmin graph voting.go             # analyze Go file
```

---

## Commands

```bash
toulmin graph <yaml>                          # validate + generate
toulmin graph <yaml> --dry-run                # print to stdout
toulmin graph <yaml> --output path/out.go     # custom output
toulmin graph <yaml> --package mypkg          # override package
toulmin graph <yaml> --check                  # validate only
toulmin graph <file.go>                       # analyze defeat cycles
toulmin evaluate                              # run example
```

---

## Evaluation Flow

```
0. Cycle detection (DFS) → error if cycle found
1. Each rule node → run func(ctx, backing) → false? skip
2. If true → traverse attackers recursively
3. Each attacker: func → false? contributes 0 → true? recurse deeper
4. h-Categoriser: raw(a) = w(a) / (1 + Σ raw(attackers)), verdict = 2*raw - 1
5. Cache per ruleID — each rule runs at most once
6. Only reachable rules executed. Backing passed as 3rd arg.
```

---

## Common Mistakes

| Mistake | Fix |
|---|---|
| Rule func wrong signature | `func(ctx Context, backing Backing) (bool, any)` |
| Chaining calls | Rule/Counter/Except return `*Rule` with builder pattern — backing/qualifier are optional |
| Attacks without registration | Must Rule/Counter/Except first to get `*Rule` |
| Verdict 0.0 as allow/deny | 0.0 = undecided — threshold is framework's decision |
| Confusing context and backing | ctx = per-request facts via Get/Set, backing = fixed criteria at declaration |
| Forgetting backing | Use `nil` when no backing needed |
| Func field in Backing struct | `Validate()` rejects func fields — use plain data fields only |

### Backing Replaces Closures

Same function + different backing values — no closure factories needed. `ruleID` distinguishes them.

```go
g := toulmin.NewGraph("example")
r1 := g.Counter(HasRole).Backing(&RoleBacking{Role: "admin"})   // ruleID = "HasRole#admin"
r2 := g.Counter(HasRole).Backing(&RoleBacking{Role: "editor"})  // ruleID = "HasRole#editor"
r1.Attacks(someRule)
```

### Verdict 0.0 Threshold

- **Security**: `verdict <= 0` → deny
- **Moderation**: `verdict < 0.3` → flag
- **Feature flags**: `verdict > 0` → enabled

Engine computes verdict. Framework interprets it.
