# toulmin

**Stop nesting if-else. Declare rules, declare relationships.**

A rule engine for Go. Rules are Go functions. Exceptions are graph edges. No DSL. No sidecar. No new language.

Rules are Go functions. Each stays at 1-2 depth:

```go
func isAuthenticated(ctx Context, backing Backing) (bool, any) {
    req, _ := ctx.Get("req")
    return req.(*Req).User != nil, nil
}
func isIPBlocked(ctx Context, backing Backing) (bool, any) {
    req, _ := ctx.Get("req")
    return blockedIPs[req.(*Req).IP], nil
}
func isInternalIP(ctx Context, backing Backing) (bool, any) {
    req, _ := ctx.Get("req")
    return strings.HasPrefix(req.(*Req).IP, "10."), nil
}
func isRateLimited(ctx Context, backing Backing) (bool, any) { /* ... */ }
func isPremiumUser(ctx Context, backing Backing) (bool, any) { /* ... */ }
func isIncidentMode(ctx Context, backing Backing) (bool, any) { /* ... */ }
```

Requirements evolve. Watch how each side handles it:

```go
// Monday: "authenticated users can access, but block banned IPs,
//          except internal network is exempt from IP blocking"
g := toulmin.NewGraph("api:access")
auth    := g.Rule(isAuthenticated)
blocked := g.Counter(isIPBlocked)
exempt  := g.Except(isInternalIP)
blocked.Attacks(auth)
exempt.Attacks(blocked)

// Tuesday: "add rate limiting"
limited := g.Counter(isRateLimited)
limited.Attacks(auth)

// Wednesday: "premium users bypass rate limit"
premium := g.Except(isPremiumUser)
premium.Attacks(limited)

// Thursday: "but not during incident response"
incident := g.Counter(isIncidentMode)
incident.Attacks(premium)

ctx := toulmin.NewContext()
ctx.Set("req", req)
results, _ := g.Evaluate(ctx)
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
func(ctx Context, backing Backing) (bool, any)
```

- `ctx` = context with Get/Set for per-request facts (user, IP, context)
- `backing` = fixed criteria at graph declaration (threshold, role name, config)
- Returns `(judgment, evidence)`. Evidence is any domain-specific type.

```go
func isInRole(ctx Context, backing Backing) (bool, any) {
    user, _ := ctx.Get("user")
    b := backing.(*RoleBacking)
    return user.(*User).Role == b.Role, user.(*User).Role
}
```

### Defeats Graph

Three node types, one edge type:

| Node | Role |
|---|---|
| **Rule** | Claim — can be attacked |
| **Counter** | Counter-claim — attacks a rule (has own conclusion) |
| **Except** | Exception — attacks only (no own conclusion) |

```go
g := toulmin.NewGraph("voting")
auth     := g.Rule(isAdult)
criminal := g.Counter(hasCriminalRecord)
expunged := g.Except(isExpunged)
criminal.Attacks(auth)      // criminal record attacks adult
expunged.Attacks(criminal)  // expungement neutralizes criminal record
```

**Except is the differentiator.** "Exceptions to exceptions" cannot be structurally expressed with if-else. In a defeats graph, it's one edge.

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
ctx := toulmin.NewContext()
ctx.Set("req", req)
results, _ := g.Evaluate(ctx)                                                         // default (matrix)
results, _ = g.Evaluate(ctx, toulmin.EvalOption{Method: toulmin.Recursive})            // recursive h-Categoriser
results, _ = g.Evaluate(ctx, toulmin.EvalOption{Trace: true})                          // with trace
results, _ = g.Evaluate(ctx, toulmin.EvalOption{Duration: true})                       // with duration (trace auto-enabled)
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
admin  := g.Rule(isInRole).Backing(&RoleBacking{Role: "admin"})
editor := g.Rule(isInRole).Backing(&RoleBacking{Role: "editor"}).Qualifier(0.8)
```

`nil` backing means the rule needs no judgment criteria. Func fields in Backing structs are forbidden — `Validate()` rejects them.

## Trace

`EvalOption{Trace: true}` tracks each rule's judgment basis:

```go
ctx := toulmin.NewContext()
results, _ := g.Evaluate(ctx, toulmin.EvalOption{Trace: true})
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
ctx := toulmin.NewContext()
ctx.Set("req", req)
results, _ := g.Evaluate(ctx)
```

`ParseYAML` parses YAML into `GraphDef`. `ValidateGraphDef` checks defeat edge references and cycles. `LoadGraph` builds the live graph from any `GraphDef` — YAML, DB, or API.

Compiled execution speed + dynamic rule updates. No DSL parser, no interpreter, no VM — just graph rewiring.

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

Example — the same access control graph from above:

```yaml
graph: api:access
rules:
  - name: isAuthenticated
    role: rule
  - name: isIPBlocked
    role: counter
  - name: isInternalIP
    role: except
  - name: isRateLimited
    role: counter
  - name: isPremiumUser
    role: except
  - name: isIncidentMode
    role: counter
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

- Adding a rule: `new.Attacks(existing)` one line vs refactoring entire nesting
- Exception handling: edge declaration vs conditions inside conditions
- Audit trail: `Evaluate(EvalOption{Trace: true})` built-in vs separate logging
- Testing: unit test each rule function vs combinatorial explosion

### vs OPA/Casbin/Cedar

| | toulmin | OPA | Casbin | Cedar |
|---|---|---|---|---|
| Rule language | **Go functions** | Rego (DSL) | PERM model (config) | Cedar (DSL) |
| Exception handling | **defeats graph** | rule priority | policy priority | forbid/permit |
| Exception of exception | **Except** | none | none | none |
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
