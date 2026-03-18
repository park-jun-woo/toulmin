# toulmin — Manual for AI Agents

Toulmin argumentation-based rule engine for Go. Rules are Go functions. Engine builds defeats graph and computes verdict via h-Categoriser.

## How to Navigate

1. Read `codebook.yaml` — project vocabulary (required/optional keys and allowed values)
2. `rg '//tm:backing'` — find why each rule exists
4. Full read only the files you need, then work

---

## Core Concepts

### Rule

All rules share one signature:

```go
func(claim any, ground any) bool
```

A rule returns `true` when its condition holds. The engine doesn't look inside — it only needs the bool result.

### Strength

| Strength | Effect |
|---|---|
| Strict | Rejects all incoming attack edges |
| Defeasible | Accepts incoming attack edges |
| Defeater | Outgoing attack edges only, no own verdict |

### Defeats Graph

Rules form a directed graph. Edges represent attacks:

- **Warrant**: node that can be attacked (primary rule)
- **Rebuttal**: node that attacks a warrant (has own conclusion)
- **Defeater**: node that attacks without own conclusion (blocks only)

The distinction is not in the type — it's in the **graph position** (defeat edges).

### Verdict

h-Categoriser computes verdict on [-1, +1] scale:

```
raw(a) = w(a) / (1 + Σ raw(attackers))     [0, 1]
verdict(a) = 2 × raw(a) - 1                [-1, 1]
```

| Verdict | Meaning |
|---|---|
| +1.0 | Violation confirmed (warrant holds, no attacks) |
| 0.0 | Undecided (warrant fully neutralized) |
| -1.0 | Fully rebutted |
| > 0 | Warrant prevails |
| < 0 | Rebuttal prevails |

### Qualifier

`0.0–1.0` float. Initial weight for h-Categoriser per rule. Default: `1.0` (Toulmin original model).

---

## Writing Rules

### Annotations

Write `//tm:backing` comment directly above the function declaration. Everything else (role, qualifier, defeats, strength) is declared at graph level.

```go
//tm:backing "Böhm-Jacopini theorem"
func CheckOneFileOneFunc(claim any, ground any) bool {
    return ground.(FileGround).FuncCount > 1
}

//tm:backing "test files conventionally group multiple test funcs"
func TestFileException(claim any, ground any) bool {
    return strings.HasSuffix(claim.(string), "_test.go")
}
```

| Annotation | Required | Description |
|---|---|---|
| `//tm:backing` | optional | Why this rule exists. Quoted string |

---

## Using the Engine (Library)

### Graph Builder API (recommended)

Functions are identifiers — no string names needed. Defeats are declared on the graph, not the function. Same function can be reused across different graphs with different defeats.

```go
g := toulmin.NewGraph("voting").
    Warrant(IsAdult, 1.0).
    Warrant(IsCitizen, 1.0).
    Rebuttal(HasCriminalRecord, 1.0).
    Defeat(HasCriminalRecord, IsAdult)

// Evaluate — verdict only (lightweight)
results := g.Evaluate(claim, ground)
for _, r := range results {
    fmt.Printf("%s: verdict=%f\n", r.Name, r.Verdict)
}

// EvaluateTrace — verdict + trace (for explainability)
results = g.EvaluateTrace(claim, ground)
for _, r := range results {
    fmt.Printf("%s: verdict=%f\n", r.Name, r.Verdict)
    for _, t := range r.Trace {
        fmt.Printf("  %s(%s): activated=%v, qualifier=%g\n", t.Name, t.Role, t.Activated, t.Qualifier)
    }
}
```

#### Evaluate vs EvaluateTrace

| Method | Returns | Use case |
|---|---|---|
| `Evaluate` | verdict only | 판정만 필요할 때 |
| `EvaluateTrace` | verdict + trace | 사유 설명이 필요할 때 |

TraceEntry contains all rules (activated and inactive) with their role, result, and qualifier:

```go
type TraceEntry struct {
    Name      string  `json:"name"`
    Role      string  `json:"role"`      // "warrant", "rebuttal", "defeater"
    Activated bool    `json:"activated"` // func(claim, ground) result
    Qualifier float64 `json:"qualifier"` // applied weight
}
```

#### Rule Reuse

```go
// Same IsAdult function, different graphs, different defeats
votingGraph := toulmin.NewGraph("voting").
    Warrant(IsAdult, 1.0).
    Rebuttal(HasCriminalRecord, 1.0).
    Defeat(HasCriminalRecord, IsAdult)

contractGraph := toulmin.NewGraph("contract").
    Warrant(IsAdult, 1.0).
    Rebuttal(IsBankrupt, 1.0).
    Defeat(IsBankrupt, IsAdult)
```

#### Qualifier Default

```go
g := toulmin.NewGraph("example").
    Warrant(IsAdult).        // qualifier = 1.0 (Toulmin original)
    Warrant(IsCitizen, 0.7)  // qualifier = 0.7 (extended use)
```

### Engine API (Phase 001, still available)

```go
eng := toulmin.NewEngine()

eng.Register(toulmin.RuleMeta{
    Name:      "CheckOneFileOneFunc",
    Qualifier: 1.0,
    Strength:  toulmin.Strict,
    Backing:   "Böhm-Jacopini theorem",
    What:      "F1: one func per file",
    Fn:        CheckOneFileOneFunc,
})

results := eng.Evaluate(claim, ground)
```

---

## Commands

```bash
toulmin graph <yaml>                          # generate graph_gen.go from YAML
toulmin graph <yaml> --dry-run                # print to stdout
toulmin graph <yaml> --output path/out.go     # custom output path
toulmin graph <yaml> --package mypkg          # override package name
toulmin evaluate                              # run example evaluation
```

### toulmin graph

1. Reads YAML graph definition
2. Validates defeats references (unknown target → error, exit 1)
3. Generates `graph_gen.go` with Graph Builder code

### YAML Graph Definition

```yaml
graph: voting
rules:
  - name: IsAdult
    role: warrant
    qualifier: 1.0
  - name: IsCitizen
    role: warrant
    qualifier: 0.7
  - name: HasCriminalRecord
    role: rebuttal
    qualifier: 1.0
defeats:
  - from: HasCriminalRecord
    to: IsAdult
```

Generated output:

```go
// Code generated by toulmin. DO NOT EDIT.
package mypkg

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

var VotingGraph = toulmin.NewGraph("voting").
    Warrant(IsAdult, 1.0).
    Warrant(IsCitizen, 0.7).
    Rebuttal(HasCriminalRecord, 1.0).
    Defeat(HasCriminalRecord, IsAdult)
```

---

## Evaluation Flow

```
1. Start from each warrant node
2. Run warrant func(claim, ground) → false? skip
3. If true, traverse attackers (rebuttal/defeater) recursively
4. Each attacker: run func → false? contributes 0 → true? recurse deeper
5. h-Categoriser at each node: raw(a) = w(a) / (1 + Σ raw(attackers))
   verdict(a) = 2 * raw(a) - 1
   Circular attack: maxDepth(100) returns 0.0
6. Func results cached — each func runs at most once per evaluation
7. Final judgment: verdict > 0 → violation, == 0 → undecided, < 0 → rebutted
```

Only rules reachable from the warrant's attack chain are executed.
EvaluateTrace returns per-warrant trace with only relevant rules.

---

## Package Structure

```
pkg/toulmin/                — public library (engine core)
  engine.go                 — Engine struct
  engine_register.go        — Engine.Register method
  engine_evaluate.go        — Engine.Evaluate method (verdict only)
  engine_evaluate_trace.go  — Engine.EvaluateTrace method (verdict + trace)
  graph_builder.go          — GraphBuilder (NewGraph, Warrant, Rebuttal, Defeater, Defeat)
  graph_builder_evaluate.go — GraphBuilder.Evaluate method (verdict only)
  graph_builder_evaluate_trace.go — GraphBuilder.EvaluateTrace method (verdict + trace)
  trace_entry.go            — TraceEntry struct
  infer_role.go             — role inference for Engine API
  func_name.go              — function pointer → name extraction
  calc_acceptability.go     — h-Categoriser computation
  build_subgraph.go         — defeats graph construction
  rule_meta.go              — RuleMeta struct
  strength.go               — Strict/Defeasible/Defeater
  node.go                   — graph node
  rule_graph.go             — RuleGraph struct
  eval_result.go            — EvalResult struct
  parse_annotation.go       — //rule: parser

internal/graphdef/          — YAML graph definition parser
internal/scanner/           — Go AST source scanner
internal/graph/             — defeats graph validation
internal/codegen/           — code generation (Graph Builder + RegisterAll)
internal/cli/               — cobra commands

cmd/toulmin/                — CLI entrypoint
```

---

## Verdict Examples

### Warrant only (no attacks)

```
W "one-func-per-file" w=1.0, attackers: none
raw = 1.0 / (1+0) = 1.0 → verdict = +1.0 (violation confirmed)
```

### Warrant + defeater

```
W "one-func-per-file" w=1.0
  attacker: D "test-file-exception" w=1.0

raw(D) = 1.0 → raw(W) = 1.0/(1+1.0) = 0.5 → verdict(W) = 0.0 (undecided)
```

### Compensation (defense)

```
W "one-func-per-file" w=1.0
  attacker: D1 "test-file-exception" w=1.0
    attacker: D2 "test-helper-not-exception" w=1.0

raw(D2) = 1.0
raw(D1) = 1.0/(1+1.0) = 0.5
raw(W)  = 1.0/(1+0.5) = 0.667 → verdict(W) = +0.33 (partially restored)
```

### Strict warrant (rejects attack)

```
W "one-func-per-file" w=1.0 strength=strict
  attacker: D — edge rejected

verdict(W) = +1.0 (unchanged)
```

---

## Common Mistakes

| Mistake | Fix |
|---|---|
| Missing `//tm:backing` | Optional, but recommended — document why the rule exists |
| Declaring role/defeats on function | Role, defeats, qualifier, strength belong in Graph Builder or YAML |
| Rule func wrong signature | Must be `func(claim any, ground any) bool` |
| Editing `graph_gen.go` manually | Re-run `toulmin graph` instead. File is generated |
