# toulmin

**Stop nesting if-else. Declare rules, declare relationships.**

A rule engine for Go. Rules are Go functions. Exceptions are graph edges. No DSL. No sidecar. No new language.

```go
// Rules are plain Go functions
func isAuthenticated(claim, ground, backing any) (bool, any) {
    return ground.(*Request).User != nil, nil
}
func isIPBlocked(claim, ground, backing any) (bool, any) {
    ip := ground.(*Request).IP
    return blockedIPs[ip], ip
}
func isInternalIP(claim, ground, backing any) (bool, any) {
    return strings.HasPrefix(ground.(*Request).IP, "10."), nil
}

// Declare relationships — don't nest conditions
g := toulmin.NewGraph("api:access")
auth    := g.Warrant(isAuthenticated, nil, 1.0)    // claim: authenticated
blocked := g.Rebuttal(isIPBlocked, nil, 1.0)       // rebuttal: IP blocked
exempt  := g.Defeater(isInternalIP, nil, 1.0)      // exception: internal network
g.Defeat(blocked, auth)    // blocked attacks auth
g.Defeat(exempt, blocked)  // internal network neutralizes block

results, _ := g.Evaluate(nil, &Request{User: user, IP: ip})
// results[0].Verdict > 0: allow
```

The same logic with if-else:

```go
if user != nil {
    if blockedIPs[ip] {
        if strings.HasPrefix(ip, "10.") {
            allow = true  // exception buried deep in nesting
        }
    } else {
        allow = true
    }
}
// Where do you insert a new rule?
// What about exceptions to exceptions? 3-4 levels deep?
```

With toulmin, it's one line: `g.Defeat(newRule, existingRule)`. Nesting depth stays flat.

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
    role := backing.(string)
    return user.Role == role, user.Role
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

### Backing

Same function + different backing = different rule. Reuse without closure factories:

```go
g := toulmin.NewGraph("access")
admin  := g.Warrant(isInRole, "admin", 1.0)
editor := g.Warrant(isInRole, "editor", 0.8)
```

`nil` backing means the rule needs no judgment criteria.

## Trace

`EvaluateTrace` tracks each rule's judgment basis:

```go
results, _ := g.EvaluateTrace(claim, ground)
for _, t := range results[0].Trace {
    fmt.Printf("%s role=%s activated=%v evidence=%v\n",
        t.Name, t.Role, t.Activated, t.Evidence)
}
```

Moderation logs, audit trails, debugging — built into the engine, no extra logging.

## Framework Packages

Domain-specific frameworks built on the core. Pre-built rule functions and wrappers.

| Package | Domain | Key API |
|---|---|---|
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
- Audit trail: `EvaluateTrace` built-in vs separate logging
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
