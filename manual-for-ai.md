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
func(claim any, ground any, backing any) (bool, any)
```

A rule returns `(true/false, evidence)`. The bool is the judgment, the second value is domain-specific evidence (e.g. error details). Return `nil` when no evidence is needed. The `backing` parameter receives the judgment criteria registered at graph declaration time.

Legacy signature `func(claim any, ground any) (bool, any)` is still supported — the engine wraps it internally via `toRuleFunc`.

### Ground vs Backing

| | Ground | Backing |
|---|---|---|
| What | Facts about the judgment target | Judgment criteria |
| When | Changes per request (runtime) | Fixed at graph declaration time |
| Passed by | `Evaluate(claim, ground)` caller | `Warrant(fn, backing, qualifier)` declaration |
| Example | File AST, user profile, request context | Threshold value, policy name, config |

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

### Rule Identity (ruleID)

Each registered rule is identified by `ruleID = funcID + "#" + backing` (when backing is non-nil). When backing is nil, `ruleID = funcID`. This allows the same function to be registered multiple times with different backing values, each as a distinct rule.

---

## Writing Rules

### Annotations

Write `//tm:backing` comment directly above the function declaration to document why the rule exists. The annotation is for documentation — the actual backing value is passed at graph declaration time via the API.

```go
//tm:backing "Böhm-Jacopini theorem"
func CheckOneFileOneFunc(claim any, ground any, backing any) (bool, any) {
    gf := ground.(*FileGround)
    if len(gf.Funcs) > 1 {
        return true, &Evidence{Got: len(gf.Funcs), Expected: 1}
    }
    return false, nil
}

//tm:backing "test files conventionally group multiple test funcs"
func TestFileException(claim any, ground any, backing any) (bool, any) {
    return strings.HasSuffix(claim.(string), "_test.go"), nil
}
```

| Annotation | Required | Description |
|---|---|---|
| `//tm:backing` | optional | Why this rule exists. Quoted string (documentation only) |

---

## Using the Engine (Library)

### Graph Builder API (recommended)

Functions are identifiers — no string names needed. Defeats are declared on the graph, not the function. Same function can be reused across different graphs with different defeats. The API takes `fn any, backing any, qualifier float64` in that order.

```go
g := toulmin.NewGraph("voting").
    Warrant(IsAdult, nil, 1.0).
    Warrant(IsCitizen, nil, 1.0).
    Rebuttal(HasCriminalRecord, nil, 1.0).
    Defeat(HasCriminalRecord, IsAdult)

// Evaluate — verdict only (lightweight)
results, err := g.Evaluate(claim, ground)
// err != nil if defeat graph contains a cycle
for _, r := range results {
    fmt.Printf("%s: verdict=%f\n", r.Name, r.Verdict)
}

// EvaluateTrace — verdict + trace (for explainability)
results, err = g.EvaluateTrace(claim, ground)
for _, r := range results {
    fmt.Printf("%s: verdict=%f\n", r.Name, r.Verdict)
    for _, t := range r.Trace {
        fmt.Printf("  %s(%s): activated=%v, qualifier=%g, backing=%v\n", t.Name, t.Role, t.Activated, t.Qualifier, t.Backing)
    }
}
```

#### Evaluate vs EvaluateTrace

| Method | Returns | Use case |
|---|---|---|
| `Evaluate` | ([]EvalResult, error) | 판정 + 증거. 순환 시 error |
| `EvaluateTrace` | ([]EvalResult, error) | 판정 + 증거 + 사유 설명. 순환 시 error |

EvalResult:

```go
type EvalResult struct {
    Name     string       `json:"name"`
    Verdict  float64      `json:"verdict"`
    Evidence any          `json:"evidence,omitempty"` // warrant's evidence
    Trace    []TraceEntry `json:"trace"`
}
```

TraceEntry contains relevant rules with their role, result, qualifier, evidence, and backing:

```go
type TraceEntry struct {
    Name      string  `json:"name"`
    Role      string  `json:"role"`      // "warrant", "rebuttal", "defeater"
    Activated bool    `json:"activated"` // func(claim, ground, backing) result
    Qualifier float64 `json:"qualifier"` // applied weight
    Evidence  any     `json:"evidence,omitempty"` // rule's evidence
    Backing   any     `json:"backing,omitempty"`  // rule's backing value
}
```

#### Rule Reuse

```go
// Same IsAdult function, different graphs, different defeats
votingGraph := toulmin.NewGraph("voting").
    Warrant(IsAdult, nil, 1.0).
    Rebuttal(HasCriminalRecord, nil, 1.0).
    Defeat(HasCriminalRecord, IsAdult)

contractGraph := toulmin.NewGraph("contract").
    Warrant(IsAdult, nil, 1.0).
    Rebuttal(IsBankrupt, nil, 1.0).
    Defeat(IsBankrupt, IsAdult)
```

#### Same Function, Different Backing

The same function can be registered multiple times with different backing values. Each registration becomes a distinct rule identified by `ruleID = funcID + "#" + backing`. Use `DefeatWith` to reference rules by both function and backing.

```go
// Same CheckThreshold function, different backing values
g := toulmin.NewGraph("limits").
    Warrant(CheckThreshold, 100, 1.0).       // ruleID = "CheckThreshold#100"
    Warrant(CheckThreshold, 200, 0.8).       // ruleID = "CheckThreshold#200"
    Rebuttal(HasExemption, "vip", 1.0).
    DefeatWith(HasExemption, "vip", CheckThreshold, 100)
```

#### DefeatWith

`DefeatWith(fromFn, fromBacking, toFn, toBacking)` declares a defeat edge between rules identified by both function and backing. Use this when the same function is registered with different backing values.

```go
g := toulmin.NewGraph("example").
    Warrant(CheckLimit, "daily", 1.0).
    Warrant(CheckLimit, "monthly", 1.0).
    Rebuttal(HasOverride, "daily", 1.0).
    DefeatWith(HasOverride, "daily", CheckLimit, "daily")  // targets only the daily variant
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

results, err := eng.Evaluate(claim, ground)
```

---

## Cycle Detection

Cyclic defeat graphs (A defeats B, B defeats A) are detected at evaluation time via DFS and return an error. The CLI can also detect cycles before code generation or in existing Go files:

```bash
toulmin graph voting.yaml --check   # validate YAML for cycles (no codegen)
toulmin graph voting.go             # analyze Go file for defeat cycles
```

---

## Commands

```bash
toulmin graph <yaml>                          # validate + generate graph_gen.go
toulmin graph <yaml> --dry-run                # print to stdout
toulmin graph <yaml> --output path/out.go     # custom output path
toulmin graph <yaml> --package mypkg          # override package name
toulmin graph <yaml> --check                  # validate only (no codegen)
toulmin graph <file.go>                       # analyze Go file for defeat cycles
toulmin evaluate                              # run example evaluation
```

### toulmin graph

1. Reads YAML graph definition (or Go source file)
2. Validates defeats references (unknown target → error, exit 1)
3. Detects cycles in defeat graph (cycle → error, exit 1)
4. Generates `graph_gen.go` with Graph Builder code (YAML only, skipped with --check)

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
    Warrant(IsAdult, nil, 1.0).
    Warrant(IsCitizen, nil, 0.7).
    Rebuttal(HasCriminalRecord, nil, 1.0).
    Defeat(HasCriminalRecord, IsAdult)
```

---

## Evaluation Flow

```
0. Cycle detection: DFS on defeat edges → error if cycle found (before any func execution)
1. Start from each warrant node
2. Run warrant func(claim, ground, backing) → false? skip
3. If true, traverse attackers (rebuttal/defeater) recursively
4. Each attacker: run func(claim, ground, backing) → false? contributes 0 → true? recurse deeper
5. h-Categoriser at each node: raw(a) = w(a) / (1 + Σ raw(attackers))
   verdict(a) = 2 * raw(a) - 1
6. Func results cached — each rule (ruleID) runs at most once per evaluation
7. Final judgment: verdict > 0 → violation, == 0 → undecided, < 0 → rebutted
```

Only rules reachable from the warrant's attack chain are executed.
EvaluateTrace returns per-warrant trace with only relevant rules.
Each rule's registered backing value is passed as the third argument to the rule function.

---

## Package Structure

```
pkg/toulmin/                — public library (engine core)
  engine.go                 — Engine struct
  engine_register.go        — Engine.Register method
  engine_evaluate.go        — Engine.Evaluate method (verdict only)
  engine_evaluate_trace.go  — Engine.EvaluateTrace method (verdict + trace)
  graph_builder.go          — GraphBuilder (NewGraph, Warrant, Rebuttal, Defeater, Defeat, DefeatWith)
  graph_builder_evaluate.go — GraphBuilder.Evaluate method (verdict only)
  graph_builder_evaluate_trace.go — GraphBuilder.EvaluateTrace method (verdict + trace)
  graph_builder_defeat_with.go — DefeatWith method (defeat with explicit backing)
  trace_entry.go            — TraceEntry struct (includes Backing field)
  infer_role.go             — role inference for Engine API
  func_name.go              — FuncName(any) → name extraction
  func_id.go                — funcID(any) → unique function identifier
  rule_id.go                — ruleID(fn, backing) → funcID + "#" + backing
  to_rule_func.go           — toRuleFunc: converts fn any to 3-arg rule func
  wrap_legacy.go            — wrapLegacy: wraps 2-arg func as 3-arg
  detect_cycle.go           — DFS cycle detection on defeat graph (exported)
  eval_context.go           — evalContext shared state
  new_eval_context.go       — evalContext factory
  eval_context_calc.go      — h-Categoriser lazy calc (single source of truth)
  eval_context_calc_trace.go — h-Categoriser lazy calc with trace collection
  eval_context_reset.go     — per-warrant state reset
  is_warrant.go             — warrant identification helper
  rule_meta.go              — RuleMeta struct (Fn takes 3-arg signature)
  strength.go               — Strict/Defeasible/Defeater
  eval_result.go            — EvalResult struct
  parse_annotation.go       — //tm: annotation parser

internal/graphdef/          — YAML graph definition parser + cycle validation
internal/analyzer/          — Go AST analysis for GraphBuilder defeat extraction
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
| Rule func wrong signature | Must be `func(claim any, ground any, backing any) (bool, any)` (legacy 2-arg also accepted) |
| Editing `graph_gen.go` manually | Re-run `toulmin graph` instead. File is generated |
| Using fn in Defeat without registering it | Must register via Warrant/Rebuttal/Defeater first |
| Treating verdict 0.0 as allow or deny | 0.0 is undecided — threshold is your framework's decision |
| Confusing ground and backing | ground = per-request facts, backing = fixed judgment criteria at graph declaration |
| Forgetting backing parameter | All Warrant/Rebuttal/Defeater calls require `(fn, backing, qualifier)`. Use `nil` when no backing is needed |

### Backing Replaces Closures

Previously, parameterized rules required closure factories (e.g. `HasRole("admin")`) which caused `funcID` identity issues because each call created a new function pointer. With backing as a first-class parameter, the same function can be registered with different backing values — no closures needed.

```go
// OLD (closure approach — no longer needed)
hasAdmin := HasRole("admin")      // closure factory
hasEditor := HasRole("editor")    // different function pointer
g := toulmin.NewGraph("example").
    Rebuttal(hasAdmin, nil, 1.0).
    Rebuttal(hasEditor, nil, 1.0).
    Defeat(hasAdmin, SomeWarrant)

// NEW (backing approach — recommended)
func HasRole(claim any, ground any, backing any) (bool, any) {
    role := backing.(string)  // backing carries the parameter
    user := ground.(*User)
    return user.HasRole(role), nil
}

g := toulmin.NewGraph("example").
    Rebuttal(HasRole, "admin", 1.0).       // ruleID = "HasRole#admin"
    Rebuttal(HasRole, "editor", 1.0).      // ruleID = "HasRole#editor"
    DefeatWith(HasRole, "admin", SomeWarrant, nil)
```

### Defeat Requires Registration

A function referenced in `Defeat(fn, target)` must be registered in the graph via `Warrant`, `Rebuttal`, or `Defeater`. Unregistered functions are not in `fnMap` and will not execute.

```go
// WRONG — whitelisted is not registered
g := toulmin.NewGraph("example").
    Warrant(IsAuthenticated, nil, 1.0).
    Rebuttal(IsIPBlocked, nil, 1.0).
    Defeat(IsIPBlocked, IsAuthenticated).
    Defeat(IsWhitelisted, IsIPBlocked)          // IsWhitelisted not in fnMap!

// CORRECT — register as Defeater
g := toulmin.NewGraph("example").
    Warrant(IsAuthenticated, nil, 1.0).
    Rebuttal(IsIPBlocked, nil, 1.0).
    Defeater(IsWhitelisted, nil, 1.0).
    Defeat(IsIPBlocked, IsAuthenticated).
    Defeat(IsWhitelisted, IsIPBlocked)
```

### Verdict 0.0 Threshold Is Framework's Decision

Verdict `0.0` means "undecided" — the warrant is exactly neutralized by its attackers. Whether `0.0` means allow or deny depends on the domain:

- **Security/route guard**: `verdict <= 0` → deny (safe default)
- **Content moderation**: `verdict < 0.3` → flag for manual review
- **Feature flags**: `verdict > 0` → enabled

The toulmin engine computes the verdict. The framework layer interprets it.
