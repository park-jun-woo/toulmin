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
func(claim any, ground any, backing any) (bool, any)
```

Returns `(judgment, evidence)`. `backing` receives judgment criteria from graph declaration. Legacy `func(any, any) (bool, any)` supported via internal wrapping.

### Ground vs Backing

| | Ground | Backing |
|---|---|---|
| What | Facts about judgment target | Judgment criteria |
| When | Per request (runtime) | Fixed at declaration |
| Passed by | `Evaluate(claim, ground)` | `Warrant(fn, backing, qualifier)` |
| Example | User, IP, request context | Threshold, role name, config |

### Strength

| Strength | Effect |
|---|---|
| Strict | Rejects all incoming attack edges |
| Defeasible | Accepts incoming attack edges |
| Defeater | Outgoing attack edges only, no own verdict |

### Defeats Graph

- **Warrant**: node that can be attacked
- **Rebuttal**: node that attacks a warrant (has own conclusion)
- **Defeater**: node that attacks without own conclusion

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
func CheckOneFileOneFunc(claim any, ground any, backing any) (bool, any) {
    gf := ground.(*FileGround)
    if len(gf.Funcs) > 1 {
        return true, &Evidence{Got: len(gf.Funcs), Expected: 1}
    }
    return false, nil
}
```

---

## Graph API

Warrant/Rebuttal/Defeater return `*Rule` reference. Defeat takes two `*Rule`. No chaining.

```go
g := toulmin.NewGraph("voting")
w := g.Warrant(IsAdult, nil, 1.0)
r := g.Rebuttal(HasCriminalRecord, nil, 1.0)
g.Defeat(r, w)

results, err := g.Evaluate(claim, ground)       // verdict only
results, err = g.EvaluateTrace(claim, ground)    // verdict + trace
```

### EvalResult / TraceEntry

```go
type EvalResult struct {
    Name     string       `json:"name"`
    Verdict  float64      `json:"verdict"`
    Evidence any          `json:"evidence,omitempty"`
    Trace    []TraceEntry `json:"trace"`
}

type TraceEntry struct {
    Name      string  `json:"name"`
    Role      string  `json:"role"`
    Activated bool    `json:"activated"`
    Qualifier float64 `json:"qualifier"`
    Evidence  any     `json:"evidence,omitempty"`
    Backing   any     `json:"backing,omitempty"`
}
```

### Same Function, Different Backing

```go
g := toulmin.NewGraph("limits")
w1 := g.Warrant(CheckThreshold, 100, 1.0)
w2 := g.Warrant(CheckThreshold, 200, 0.8)
r := g.Rebuttal(HasExemption, "vip", 1.0)
g.Defeat(r, w1)
```

### LoadGraph — Dynamic Graph from Definition

Builds a live `*Graph` from a `GraphDef`, function registry, and optional backing registry. No code generation, no recompilation.

```go
funcs := map[string]any{
    "isAuthenticated": isAuthenticated,
    "isIPBlocked":     isIPBlocked,
}
backings := map[string]any{
    "isIPBlocked": fetchBlocklistFromRedis(),
}

g, err := toulmin.LoadGraph(def, funcs, backings)
results, _ := g.Evaluate(claim, ground)
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
    role: <role>           # warrant | rebuttal | defeater (required)
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
    role: warrant
  - name: isIPBlocked
    role: rebuttal
  - name: isInternalIP
    role: defeater
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
results, _ := g.Evaluate(nil, req)
```

### Engine API (legacy)

`Engine.Register(RuleMeta{...})` + `Engine.Evaluate()` — string-based names, still available.

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
1. Each warrant node → run func(claim, ground, backing) → false? skip
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
| Rule func wrong signature | `func(claim, ground, backing any) (bool, any)` (legacy 2-arg also OK) |
| Chaining calls | Warrant/Rebuttal/Defeater return `*Rule`, not `*Graph` — use separate statements |
| Defeat without registration | Must Warrant/Rebuttal/Defeater first to get `*Rule` |
| Verdict 0.0 as allow/deny | 0.0 = undecided — threshold is framework's decision |
| Confusing ground and backing | ground = per-request facts, backing = fixed criteria at declaration |
| Forgetting backing | Use `nil` when no backing needed |

### Backing Replaces Closures

Same function + different backing values — no closure factories needed. `ruleID` distinguishes them.

```go
g := toulmin.NewGraph("example")
r1 := g.Rebuttal(HasRole, "admin", 1.0)   // ruleID = "HasRole#admin"
r2 := g.Rebuttal(HasRole, "editor", 1.0)  // ruleID = "HasRole#editor"
g.Defeat(r1, someWarrant)
```

### Verdict 0.0 Threshold

- **Security**: `verdict <= 0` → deny
- **Moderation**: `verdict < 0.3` → flag
- **Feature flags**: `verdict > 0` → enabled

Engine computes verdict. Framework interprets it.
