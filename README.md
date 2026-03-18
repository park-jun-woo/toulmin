# toulmin

Toulmin argumentation-based rule engine for Go.

Rules are Go functions. The engine builds a defeats graph from activated rules and computes verdicts via Amgoud's h-Categoriser on a [-1, +1] scale.

## Install

```bash
go get github.com/park-jun-woo/toulmin/pkg/toulmin
```

## Usage

### Graph Builder API

Functions are identifiers. No string names needed.

```go
g := toulmin.NewGraph("voting").
    Warrant(IsAdult, 1.0).
    Warrant(IsCitizen, 1.0).
    Rebuttal(HasCriminalRecord, 1.0).
    Defeat(IsAdult, HasCriminalRecord)

// Evaluate — verdict + evidence
results, err := g.Evaluate(claim, ground)
// err: non-nil if defeat graph contains a cycle
// results[0].Verdict: +1.0 violation, 0.0 undecided, -1.0 rebutted
// results[0].Evidence: warrant's domain-specific evidence (any)

// EvaluateTrace — verdict + evidence + per-warrant trace
traced, err := g.EvaluateTrace(claim, ground)
// traced[0].Trace: relevant rules with name, role, activated, qualifier, evidence
```

### Rule Reuse

Same function, different graphs, different defeats:

```go
votingGraph := toulmin.NewGraph("voting").
    Warrant(IsAdult, 1.0).
    Rebuttal(HasCriminalRecord, 1.0).
    Defeat(HasCriminalRecord, IsAdult)

contractGraph := toulmin.NewGraph("contract").
    Warrant(IsAdult, 1.0).
    Rebuttal(IsBankrupt, 1.0).
    Defeat(IsBankrupt, IsAdult)
```

### Qualifier

Defaults to 1.0 (Toulmin original model). Override per rule:

```go
g := toulmin.NewGraph("example").
    Warrant(IsAdult).        // qualifier = 1.0
    Warrant(IsCitizen, 0.7)  // qualifier = 0.7
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

## Verdict

h-Categoriser: `raw(a) = w(a) / (1 + Sum(raw(attackers)))`, then `verdict = 2*raw - 1`.

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
func(claim any, ground any) (bool, any)
```

Returns `(judgment, evidence)`. Evidence is domain-specific (`any`). Return `nil` when not needed.

## Annotations

Backing stays on the function as optional metadata:

```go
//tm:backing "Bohm-Jacopini theorem"
func CheckOneFileOneFunc(claim any, ground any) (bool, any) { ... }
```

## Theory

| Component | Source |
|---|---|
| 6-element structure | Toulmin (1958) |
| strict/defeasible/defeater | Nute (1994) |
| h-Categoriser | Amgoud & Ben-Naim (2013, 2017) |

## License

MIT
