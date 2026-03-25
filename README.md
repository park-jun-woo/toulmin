# toulmin

**Stop nesting if-else. Declare rules, declare relationships.**

A rule engine for Go. Rules are Go functions. Exceptions are graph edges. No DSL. No sidecar. No new language.

Rules are Go functions. Each stays at 1-2 depth:

```go
func isAuthenticated(claim, ground, backing any) (bool, any) {
    return ground.(*Req).User != nil, nil
}
func isIPBlocked(claim, ground, backing any) (bool, any) {
    return blockedIPs[ground.(*Req).IP], nil
}
func isInternalIP(claim, ground, backing any) (bool, any) {
    return strings.HasPrefix(ground.(*Req).IP, "10."), nil
}
func isRateLimited(claim, ground, backing any) (bool, any) { /* ... */ }
func isPremiumUser(claim, ground, backing any) (bool, any) { /* ... */ }
func isIncidentMode(claim, ground, backing any) (bool, any) { /* ... */ }
```

Requirements evolve. Watch how each side handles it:

```go
// Monday: "authenticated users can access, but block banned IPs,
//          except internal network is exempt from IP blocking"
g := toulmin.NewGraph("api:access")
auth    := g.Warrant(isAuthenticated, nil, 1.0)
blocked := g.Rebuttal(isIPBlocked, nil, 1.0)
exempt  := g.Defeater(isInternalIP, nil, 1.0)
g.Defeat(blocked, auth)
g.Defeat(exempt, blocked)

// Tuesday: "add rate limiting"
limited := g.Rebuttal(isRateLimited, nil, 1.0)
g.Defeat(limited, auth)

// Wednesday: "premium users bypass rate limit"
premium := g.Defeater(isPremiumUser, nil, 1.0)
g.Defeat(premium, limited)

// Thursday: "but not during incident response"
incident := g.Rebuttal(isIncidentMode, nil, 1.0)
g.Defeat(incident, premium)

results, _ := g.Evaluate(nil, req)
// results[0].Verdict > 0: allow
```

Each day: 2 lines added, nothing else changes. Now the same evolution with if-else:

```go
// Monday
if user != nil {
    if blockedIPs[ip] {
        if strings.HasPrefix(ip, "10.") {
            allow = true
        }
    } else {
        allow = true
    }
}

// Tuesday: "add rate limiting" — where does it go?
if user != nil {
    if blockedIPs[ip] {
        if strings.HasPrefix(ip, "10.") {
            allow = true
        }
    } else if isRateLimited(ip) {
        allow = false
    } else {
        allow = true
    }
}

// Wednesday: "premium users bypass rate limit"
if user != nil {
    if blockedIPs[ip] {
        if strings.HasPrefix(ip, "10.") {
            allow = true
        }
    } else if isRateLimited(ip) {
        if isPremium(user) {       // 3 levels deep
            allow = true
        }
    } else {
        allow = true
    }
}

// Thursday: "but not during incident response"
if user != nil {
    if blockedIPs[ip] {
        if strings.HasPrefix(ip, "10.") {
            allow = true
        }
    } else if isRateLimited(ip) {
        if isPremium(user) {
            if !incidentMode {     // 4 levels, losing track
                allow = true
            }
        }
    } else {
        allow = true
    }
}
```

toulmin: **2 lines per requirement, structure never changes.** if-else: **restructure everything, every time.**

## Install

```bash
go get github.com/park-jun-woo/toulmin/pkg/toulmin
```

## Core Concepts

### Rules are Go functions

```go
func(claim any, ground any, backing any) (bool, any)
```

- `ground` = per-request facts (user, IP, context)
- `backing` = fixed criteria at graph declaration (threshold, role name, config)
- Returns `(judgment, evidence)`. Evidence is any domain-specific type.

```go
func isInRole(claim, ground, backing any) (bool, any) {
    user := ground.(*User)
    b := backing.(*RoleBacking)
    return user.Role == b.Role, user.Role
}
```

### Defeats Graph

Three node types, one edge type:

| Node | Role |
|---|---|
| **Warrant** | Claim — can be attacked |
| **Rebuttal** | Counter-claim — attacks a warrant (has own conclusion) |
| **Defeater** | Exception — attacks only (no own conclusion) |

```go
g := toulmin.NewGraph("voting")
auth     := g.Warrant(isAdult, nil, 1.0)
criminal := g.Rebuttal(hasCriminalRecord, nil, 1.0)
expunged := g.Defeater(isExpunged, nil, 1.0)
g.Defeat(criminal, auth)      // criminal record attacks adult
g.Defeat(expunged, criminal)  // expungement neutralizes criminal record
```

**Defeater is the differentiator.** "Exceptions to exceptions" cannot be structurally expressed with if-else. In a defeats graph, it's one edge.

### Verdict

h-Categoriser computes a continuous value `[-1, +1]`:

```
raw(a) = w(a) / (1 + Σ raw(attackers))
verdict = 2 × raw - 1
```

| Verdict | Meaning |
|---|---|
| +1.0 | Confirmed |
| 0.0 | Undecided |
| -1.0 | Fully rebutted |

Not binary. **The framework interprets**:

- Security: `verdict ≤ 0` → deny
- Moderation: `verdict ≤ 0` → block, `0 < v ≤ 0.3` → flag, `> 0.3` → allow
- Feature flags: `verdict > 0` → enabled

### Evaluation Options

```go
results, _ := g.Evaluate(nil, req)                                                    // default (matrix)
results, _ = g.Evaluate(nil, req, toulmin.EvalOption{Method: toulmin.Recursive})       // recursive h-Categoriser
results, _ = g.Evaluate(nil, req, toulmin.EvalOption{Trace: true})                     // with trace
results, _ = g.Evaluate(nil, req, toulmin.EvalOption{Duration: true})                  // with duration (trace auto-enabled)
```

`EvalOption` controls evaluation behavior: `Method` (Matrix/Recursive), `Trace` (collect TraceEntry), `Duration` (measure per-rule execution time).

### Backing

Backing values must implement the `Backing` interface:

```go
type Backing interface {
    BackingName() string
    Validate() error
}
```

Same function + different backing = different rule. Reuse without closure factories:

```go
g := toulmin.NewGraph("access")
admin  := g.Warrant(isInRole, &RoleBacking{Role: "admin"}, 1.0)
editor := g.Warrant(isInRole, &RoleBacking{Role: "editor"}, 0.8)
```

`nil` backing means the rule needs no judgment criteria. Func fields in Backing structs are forbidden — `Validate()` rejects them.

## Trace

`EvalOption{Trace: true}` tracks each rule's judgment basis:

```go
results, _ := g.Evaluate(claim, ground, toulmin.EvalOption{Trace: true})
for _, t := range results[0].Trace {
    fmt.Printf("%s role=%s activated=%v evidence=%v\n",
        t.Name, t.Role, t.Activated, t.Evidence)
}
```

Moderation logs, audit trails, debugging — built into the engine, no extra logging.

## Dynamic Loading

`LoadGraph` builds a live graph from a definition + function registry. Graph structure and backing change without redeployment — functions stay compiled.

```go
// Register compiled functions once at startup
funcs := map[string]any{
    "isAuthenticated": isAuthenticated,
    "isIPBlocked":     isIPBlocked,
    "isRateLimited":   isRateLimited,
}

// Parse YAML into GraphDef, validate, then build live graph
def, _ := toulmin.ParseYAML("policy.yaml")
toulmin.ValidateGraphDef(def)
backings := map[string]toulmin.Backing{"isIPBlocked": fetchBlocklistFromRedis()}
g, err := toulmin.LoadGraph(def, funcs, backings)
results, _ := g.Evaluate(nil, req)
```

`ParseYAML` parses YAML into `GraphDef`. `ValidateGraphDef` checks defeat edge references and cycles. `LoadGraph` builds the live graph from any `GraphDef` — YAML, DB, or API.

Compiled execution speed + dynamic rule updates. No DSL parser, no interpreter, no VM — just graph rewiring.

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

Example — the same access control graph from above:

```yaml
graph: api:access
rules:
  - name: isAuthenticated
    role: warrant
  - name: isIPBlocked
    role: rebuttal
  - name: isInternalIP
    role: defeater
  - name: isRateLimited
    role: rebuttal
  - name: isPremiumUser
    role: defeater
  - name: isIncidentMode
    role: rebuttal
defeats:
  - from: isIPBlocked
    to: isAuthenticated
  - from: isInternalIP
    to: isIPBlocked
  - from: isRateLimited
    to: isAuthenticated
  - from: isPremiumUser
    to: isRateLimited
  - from: isIncidentMode
    to: isPremiumUser
```

## Framework Packages

Domain-specific frameworks built on the core. Pre-built rule functions and wrappers.

| Package | Domain | Key API |
|---|---|---|
| `pkg/toulmin` | Core engine, YAML parser, validator, codegen | `Graph`, `ParseYAML`, `ValidateGraphDef`, `GenerateGraph` |
| `pkg/analyzer` | Go AST defeat graph extraction | `ExtractDefeats` |
| `pkg/policy` | Access control (auth, IP, rate limit) | `Guard` (net/http middleware) |
| `pkg/state` | State transitions (FSM) | `Machine.Can`, `Mermaid()` |
| `pkg/approve` | Multi-step approval workflow | `Flow.Evaluate` |
| `pkg/price` | Price judgment (coupons, membership) | `Pricer.Evaluate` |
| `pkg/feature` | Feature flags (rollout, toggle) | `Flags.IsEnabled` |
| `pkg/moderate` | Content moderation (hate speech, spam) | `Moderator.Review` |

You can use the core without any framework. Writing your own rule functions — like the killer example above — is the most flexible approach.

## Why toulmin

### vs if-else

- Adding a rule: `g.Defeat(new, existing)` one line vs refactoring entire nesting
- Exception handling: edge declaration vs conditions inside conditions
- Audit trail: `Evaluate(EvalOption{Trace: true})` built-in vs separate logging
- Testing: unit test each rule function vs combinatorial explosion

### vs OPA/Casbin/Cedar

| | toulmin | OPA | Casbin | Cedar |
|---|---|---|---|---|
| Rule language | **Go functions** | Rego (DSL) | PERM model (config) | Cedar (DSL) |
| Exception handling | **defeats graph** | rule priority | policy priority | forbid/permit |
| Exception of exception | **Defeater** | none | none | none |
| Judgment | **continuous [-1,1]** | allow/deny | allow/deny | allow/deny |
| Audit trail | **Trace built-in** | Decision log | none | none |
| Dependencies | Go stdlib only | Go + Rego runtime | Go | Rust + FFI |
| Learning curve | Just know Go | Learn Rego | Learn PERM model | Learn Cedar syntax |

### Academic Foundation

| Component | Source |
|---|---|
| 6-element argumentation | Toulmin (1958) |
| strict/defeasible/defeater | Nute (1994) |
| h-Categoriser | Amgoud & Ben-Naim (2013, 2017) |

## Testing

`RunCases` runs table-driven policy tests with zero boilerplate:

```go
func TestAccessPolicy(t *testing.T) {
    g := buildAccessGraph()
    toulmin.RunCases(t, g, []toulmin.TestCase{
        {Name: "admin allowed",    Ground: &Ctx{Role: "admin"},  Expect: toulmin.VerdictAbove(0)},
        {Name: "blocked IP",       Ground: &Ctx{IP: "blocked"},  Expect: toulmin.VerdictAtMost(0)},
        {Name: "unauthenticated",  Ground: &Ctx{User: nil},      Expect: toulmin.NoResult},
        {Name: "partial override", Ground: &Ctx{Role: "editor"}, Expect: toulmin.VerdictBetween(0, 0.5)},
    })
}
```

| Expectation | Condition |
|---|---|
| `VerdictAbove(v)` | verdict > v |
| `VerdictAtMost(v)` | verdict <= v |
| `VerdictBetween(lo, hi)` | lo < verdict <= hi |
| `NoResult` | no active warrants |

## CLI

```bash
toulmin graph voting.yaml                  # YAML → Go code generation
toulmin graph voting.yaml --check          # cycle validation only
toulmin graph voting.yaml --dry-run        # print to stdout
toulmin graph voting.go                    # analyze Go file for cycles
```

## Used By

- **[filefunc](https://github.com/park-jun-woo/filefunc)** — LLM-native Go code structure tool. The `validate` command uses toulmin defeats graph to handle rule exceptions (F5, F6, etc.).

## License

MIT
