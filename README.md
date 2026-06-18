# toulmin

[![version](https://img.shields.io/badge/version-0.3.0-blue)](https://github.com/park-jun-woo/toulmin/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow)](https://opensource.org/licenses/MIT)

**Stop nesting if-else. Declare rules, declare relationships.**

A rule engine for TypeScript, Python, and Go. Rules are functions. Exceptions are graph edges. No DSL. No sidecar. No new language.

### TypeScript

```typescript
const isAuthenticated = (ctx, specs) => [ctx.get("user") != null, null]
const isIPBlocked     = (ctx, specs) => [blockedIPs.has(ctx.get("ip")), null]
const isInternalIP    = (ctx, specs) => [ctx.get("ip")?.startsWith("10."), null]
const isRateLimited   = (ctx, specs) => [/* ... */]
const isPremiumUser   = (ctx, specs) => [/* ... */]
const isIncidentMode  = (ctx, specs) => [/* ... */]

const g = new Graph("api:access")
const auth    = g.rule(isAuthenticated)
const blocked = g.counter(isIPBlocked)
const exempt  = g.except(isInternalIP)
blocked.attacks(auth)
exempt.attacks(blocked)

const limited  = g.counter(isRateLimited)    // Tuesday: rate limiting
limited.attacks(auth)
const premium  = g.except(isPremiumUser)     // Wednesday: premium bypass
premium.attacks(limited)
const incident = g.counter(isIncidentMode)   // Thursday: incident override
incident.attacks(premium)

const results = g.evaluate(newContext())
// results[0].verdict > 0: allow
```

### Python (planned)

```python
g = Graph("api:access")
auth    = g.rule(is_authenticated)
blocked = g.counter(is_ip_blocked)
exempt  = g.except_(is_internal_ip)
blocked.attacks(auth)
exempt.attacks(blocked)

results = g.evaluate(MapContext())
```

### Go

```go
g := toulmin.NewGraph("api:access")
auth    := g.Rule(isAuthenticated)
blocked := g.Counter(isIPBlocked)
exempt  := g.Except(isInternalIP)
blocked.Attacks(auth)
exempt.Attacks(blocked)

results, _ := g.Evaluate(ctx)
```

Requirements evolve. Each day: 2 lines added, nothing else changes. Now the same evolution with if-else:

```go
// Monday
func isAuthenticated(ctx Context, specs Specs) (bool, any) {
    req, _ := ctx.Get("req")
    return req.(*Req).User != nil, nil
}

// Thursday — 4 levels deep
if user != nil {
    if blockedIPs[ip] {
        if strings.HasPrefix(ip, "10.") { allow = true }
    } else if isRateLimited(ip) {
        if isPremium(user) {
            if !incidentMode { allow = true }  // losing track
        }
    } else { allow = true }
}
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
npm install rulecat          # TypeScript
pip install rulecat          # Python (planned)
go get github.com/park-jun-woo/toulmin/pkg/toulmin  # Go
```

## Core Concepts

### Rules are functions

```typescript
// TypeScript
const fn: RuleFunc = (ctx, specs) => [boolean, unknown]

// Python
def fn(ctx: Context, specs: list[Spec]) -> tuple[bool, Any]: ...

// Go
func fn(ctx Context, specs Specs) (bool, any)
```

- `ctx` = context with get/set for per-request facts (user, IP, context)
- `specs` = fixed criteria at graph declaration via `.with()` / `.with_spec()` / `.With()` (threshold, role name, config)
- Returns `(judgment, evidence)`. Evidence is any domain-specific type.

```typescript
const isInRole: RuleFunc = (ctx, specs) => {
    const user = ctx.get("user")
    if (!user) return [false, null]
    const role = (specs[0] as RoleSpec).role
    return [user.role === role, user.role]
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
results, _ = g.Evaluate(ctx, toulmin.EvalOption{Trace: true})                          // with trace
results, _ = g.Evaluate(ctx, toulmin.EvalOption{Duration: true})                       // with duration (trace auto-enabled)
```

`EvalOption` controls evaluation behavior: `Method` (Matrix/Recursive (planned)), `Trace` (collect TraceEntry), `Duration` (measure per-rule execution time).

### Spec

Spec values must implement the `Spec` interface:

```go
type Spec interface {
    SpecName() string
    Validate() error
}
```

Same function + different spec = different rule. Reuse without closure factories:

```go
g := toulmin.NewGraph("access")
admin  := g.Rule(isInRole).With(&RoleSpec{Role: "admin"})
editor := g.Rule(isInRole).With(&RoleSpec{Role: "editor"}).Qualifier(0.8)
```

`nil` specs means the rule needs no judgment criteria. Func fields in Spec structs are forbidden — `Validate()` rejects them.

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

## Run

`Evaluate` judges and returns — pure and idempotent. `Run` judges first, **then acts**: it
does a full pass, then fires one handler per **Active** node. A node is Active when its rule
applied **and** prevailed (`activated && verdict > 0`) — that is the only event. To inspect
any other node's outcome, read the returned `trace`: `verdict > 0` prevailed, `verdict <= 0`
defeated, `activated == false` did not apply.

```go
g.Rule(isAuthenticated).
    RunOn(func(self toulmin.TraceEntry, t toulmin.Trace) error {
        return audit(self)                  // self = this handler's own node (ran and prevailed)
    })

results, trace, err := g.Run(ctx)  // []EvalResult, Trace, error
```

Handlers fire in registration order; the first handler error stops `Run`. The handler gets its
own node's entry as `self` plus a read-only `Trace` `t` — the same value `Run` returns — with
exactly three methods: `t.All()` (every node's `TraceEntry`, registration order), `t.Get(name)`
(one node by short name), and `t.Ctx()` (this Run's context). `self` is the first argument,
handed in by the engine, so a handler never has to look up its own row by name (which was
fragile when an entry's name carried a dynamic suffix). Each `TraceEntry` exposes **Claim** (`Name`), **Ground** (`Ground` = the `ctx`
as-is), **Backing** (`Specs`), and **Verdict** — enough to audit, explain, or apply gradient
thresholds without any separate view. There is no direct attacker lookup; reason from
`t.All()[i].Verdict`. `ctx` is mutable (side effects propagate); the trace is the judgment
record. If you serialise the trace as JSON, keep `ctx` to serialisable values (`Ground = ctx`).

### Execution composition

`rule.Run(g)` declares that when a node is **Active**, its sub-graph `g` is Run with the same
ctx — a graph of graphs. Judgment composes *upward* through `Attacks` (verdict flows up);
execution composes *downward* through `Run` (the sub-graph's verdict stays isolated, only
errors propagate). Execution composition must be a DAG — cycles are rejected up front, depth
is capped at 64.

```go
order := g.Rule(orderPlaced)
order.RunOn(logOrder).Run(notifyGraph)   // Active order → Run the notify graph
```

The `RunOn` handler + `Run(g)` execution composition is available across the Go, TypeScript, and Python ports.

## Framework Packages

Domain-specific frameworks built on the core. Pre-built rule functions and wrappers.

| Package | Domain | Key API |
|---|---|---|
| `pkg/toulmin` | Core engine | `Graph`, `EvalOption` |
| `pkg/policy` | Access control (auth, IP, rate limit) | `Guard` (net/http middleware) |
| `pkg/state` | State transitions (FSM) | `Machine.Can`, `Mermaid()` |
| `pkg/approve` | Multi-step approval workflow | `Flow.Evaluate` |
| `pkg/price` | Price judgment (coupons, membership) | `Pricer.Evaluate` |
| `pkg/feature` | Feature flags (rollout, toggle) | `Flags.IsEnabled` |
| `pkg/moderate` | Content moderation (hate speech, spam) | `Moderator.Review` |
| `pkg/tangl` | Markdown-based policy language | `parser.Parse`, `validate.Validate` |

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
| Rule language | **TS/Python/Go functions** | Rego (DSL) | PERM model (config) | Cedar (DSL) |
| Exception handling | **defeats graph** | rule priority | policy priority | forbid/permit |
| Exception of exception | **Except** | none | none | none |
| Judgment | **continuous [-1,1]** | allow/deny | allow/deny | allow/deny |
| Audit trail | **Trace built-in** | Decision log | none | none |
| Dependencies | **Zero** | Go + Rego runtime | Go | Rust + FFI |
| Learning curve | Know your language | Learn Rego | Learn PERM model | Learn Cedar syntax |

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
        {Name: "admin allowed",    Context: &Ctx{Role: "admin"},  Expect: toulmin.VerdictAbove(0)},
        {Name: "blocked IP",       Context: &Ctx{IP: "blocked"},  Expect: toulmin.VerdictAtMost(0)},
        {Name: "unauthenticated",  Context: &Ctx{User: nil},      Expect: toulmin.NoResult},
        {Name: "partial override", Context: &Ctx{Role: "editor"}, Expect: toulmin.VerdictBetween(0, 0.5)},
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
toulmin evaluate    # run example
```

## Used By

- **[filefunc](https://github.com/park-jun-woo/filefunc)** — LLM-native Go code structure tool. The `validate` command uses toulmin defeats graph to handle rule exceptions (F5, F6, etc.).

## License

MIT
