# toulmin

**Stop nesting if-else. Declare rules, declare relationships.**

A lightweight rule engine for Go. Write each rule as a function, wire exceptions as graph edges. No DSL, no sidecar, no new language to learn.

```go
// if-else — hard to read, harder to maintain
if isAdult {
    if isCitizen {
        if !hasCriminalRecord {
            if !isSuspended {
                allow = true
            }
        } else if isExpunged {
            allow = true
        }
    }
}

// toulmin — declare rules and relationships
g := toulmin.NewGraph("voting").
    Warrant(IsAdult, nil, 1.0).
    Warrant(IsCitizen, nil, 1.0).
    Rebuttal(HasCriminalRecord, nil, 1.0).
    Rebuttal(IsSuspended, nil, 1.0).
    Defeat(IsExpunged, HasCriminalRecord)
```

Each rule function stays at 1-2 depth. Complexity lives in the graph, not in nesting.

```go
func IsAdult(claim any, ground any, backing any) (bool, any) {
    user := ground.(*User)
    return user.Age >= 18, nil
}
```

## Install

```bash
go get github.com/park-jun-woo/toulmin/pkg/toulmin
```

## Usage

### Graph Builder API

Functions are identifiers. No string names needed. Backing is passed as a second argument — use `nil` when the rule needs no external criteria.

```go
g := toulmin.NewGraph("voting").
    Warrant(IsAdult, nil, 1.0).
    Warrant(IsCitizen, nil, 1.0).
    Rebuttal(HasCriminalRecord, nil, 1.0).
    Defeat(HasCriminalRecord, IsAdult)

// Evaluate — verdict + evidence
results, err := g.Evaluate(claim, ground)
// err: non-nil if defeat graph contains a cycle
// results[0].Verdict: +1.0 violation, 0.0 undecided, -1.0 rebutted
// results[0].Evidence: warrant's domain-specific evidence (any)

// EvaluateTrace — verdict + evidence + per-warrant trace
traced, err := g.EvaluateTrace(claim, ground)
// traced[0].Trace: relevant rules with name, role, activated, qualifier, evidence, backing
```

### Backing

Backing is a first-class Toulmin element — judgment criteria passed as an argument to the rule function. The same function can serve different purposes with different backing values.

```go
// Same function, different backing — "admin" vs "editor" role checks
g := toulmin.NewGraph("admin").
    Warrant(IsInRole, "admin", 1.0).
    Warrant(IsAuthenticated, nil, 1.0)

g := toulmin.NewGraph("editor").
    Warrant(IsInRole, "editor", 1.0).
    Warrant(IsAuthenticated, nil, 1.0)
```

The rule function receives backing as the third argument:

```go
func IsInRole(claim any, ground any, backing any) (bool, any) {
    user := ground.(*User)
    role := backing.(string)
    return user.Role == role, nil
}
```

When two registrations use the same function but different backing, use `DefeatWith` to specify which defeats which:

```go
g := toulmin.NewGraph("firewall").
    Warrant(IsAuthenticated, nil, 1.0).
    Rebuttal(IsIPInList, blocklist, 1.0).
    Warrant(IsIPInList, whitelist, 1.0).
    DefeatWith(IsIPInList, blocklist, IsAuthenticated, nil).
    DefeatWith(IsIPInList, whitelist, IsIPInList, blocklist)
```

### Rule Reuse

Same function, different graphs, different defeats:

```go
votingGraph := toulmin.NewGraph("voting").
    Warrant(IsAdult, nil, 1.0).
    Rebuttal(HasCriminalRecord, nil, 1.0).
    Defeat(HasCriminalRecord, IsAdult)

contractGraph := toulmin.NewGraph("contract").
    Warrant(IsAdult, nil, 1.0).
    Rebuttal(IsBankrupt, nil, 1.0).
    Defeat(IsBankrupt, IsAdult)
```

### Qualifier

Defaults to 1.0. Set to 1.0 for simple pass/fail policies. Use fractional values for weighted judgment:

```go
g := toulmin.NewGraph("example").
    Warrant(IsAdult, nil, 1.0).
    Warrant(IsCitizen, nil, 0.7)
```

### Engine API (Phase 001)

```go
eng := toulmin.NewEngine()
eng.Register(toulmin.RuleMeta{
    Name:      "OneFileOneFunc",
    Qualifier: 1.0,
    Strength:  toulmin.Defeasible,
    Fn:        CheckOneFileOneFunc,
})
results, err := eng.Evaluate(claim, ground)
```

## Verdict

Verdicts are continuous, not binary. h-Categoriser computes `raw(a) = w(a) / (1 + Sum(raw(attackers)))`, then `verdict = 2*raw - 1`.

| Verdict | Meaning |
|---|---|
| +1.0 | Violation confirmed |
| 0.0 | Undecided |
| -1.0 | Fully rebutted |

## Strength

| Strength | Effect |
|---|---|
| Strict | Rejects all incoming attack edges |
| Defeasible | Accepts incoming attack edges |
| Defeater | Outgoing attack edges only, no own verdict |

## Rule Signature

```go
func(claim any, ground any, backing any) (bool, any)
```

Returns `(judgment, evidence)`. Evidence is domain-specific (`any`). Return `nil` when not needed. Backing is the judgment criteria supplied at registration — `nil` when the rule needs no external context.

## Cycle Detection

Cyclic defeat graphs (e.g. A defeats B, B defeats A) are rejected at evaluation time with an error. The CLI can also detect cycles before code generation:

```bash
toulmin graph voting.yaml --check          # validate YAML for cycles
toulmin graph voting.go                    # analyze Go file for cycles
```

## YAML Graph Definition

Define graph structure in YAML, generate Go code:

```yaml
graph: voting
rules:
  - name: IsAdult
    role: warrant
    qualifier: 1.0
  - name: HasCriminalRecord
    role: rebuttal
    qualifier: 1.0
defeats:
  - from: HasCriminalRecord
    to: IsAdult
```

```bash
toulmin graph voting.yaml                    # validate + generate graph_gen.go
toulmin graph voting.yaml --dry-run          # print to stdout
toulmin graph voting.yaml --output out.go    # custom output path
toulmin graph voting.yaml --check            # validate only (no code generation)
toulmin graph voting.go                      # analyze Go file for defeat cycles
```

## Annotations

Backing stays on the function as optional metadata:

```go
//tm:backing "Bohm-Jacopini theorem"
func CheckOneFileOneFunc(claim any, ground any, backing any) (bool, any) { ... }
```

## Theory

| Component | Source |
|---|---|
| 6-element structure | Toulmin (1958) |
| strict/defeasible/defeater | Nute (1994) |
| h-Categoriser | Amgoud & Ben-Naim (2013, 2017) |

## Used By

- **[filefunc](https://github.com/park-jun-woo/filefunc)** — LLM-native Go code structure convention tool. The `validate` command uses toulmin defeats graph to handle rule exceptions (F5, F6, etc.).
- **pkg/route** — HTTP guard middleware built on toulmin. Composes authentication, IP blocking, rate limiting, and role checks as a defeats graph.

## License

MIT
