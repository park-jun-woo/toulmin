# pkg/toulmin

Toulmin argumentation-based rule engine core. Defeats graph + h-Categoriser verdict.

## Rule Function Signature

```go
func(ctx Context, specs Specs) (bool, any)
// ctx = Context interface with Get/Set for per-request facts
// specs = fixed criteria attached via .With() at declaration
// returns (judgment, evidence)
```

## Public API

### Graph Builder

| Function | File | Description |
|---|---|---|
| `NewGraph(name)` | new_graph.go | Create graph builder |
| `g.Rule(fn)` | graph_rule.go | Register rule, return `*Rule` (builder: `.With(spec)`, `.Qualifier(q)`) |
| `g.Counter(fn)` | graph_counter.go | Register counter, return `*Rule` (builder: `.With(spec)`, `.Qualifier(q)`) |
| `g.Except(fn)` | graph_except.go | Register except, return `*Rule` (builder: `.With(spec)`, `.Qualifier(q)`) |
| `rule.Attacks(target)` | rule_attacks.go | Declare defeat edge (method on `*Rule`) |
| `rule.With(spec)` | rule_with.go | Attach a Spec to the rule (additive, supports chaining) |
| `g.Evaluate(ctx, opt ...EvalOption)` | graph_evaluate.go | Run evaluation, return verdicts (ctx is Context, opt controls Trace/Duration) |
| `g.Run(ctx, opt ...EvalOption)` | graph_run.go | Pre-evaluate (full pass), fire each Active node's `RunOn` handler with the whole-graph `Trace`, and Run any Active node's sub-graph; returns `([]EvalResult, Trace, error)` |
| `rule.RunOn(h)` | rule_run_on.go | Register the handler fired when this node is Active (func true && verdict > 0) |
| `rule.Run(g)` | rule_run.go | Declare execution edge: when this node is Active, Run sub-graph `g` (graph-of-graphs; panics on nil) |
| `NewContext()` | new_context.go | Create `*MapContext` implementing Context interface |

`Run` forces Trace/Duration off and does a full pass, then fires the `RunOn` handler of every
**Active** node — a node that applied and prevailed (`Activated && Verdict > 0`). That is the
only event; there are no Defeated/Inactive handlers. Handlers fire in rule registration order;
a handler error or panic stops `Run` immediately and is returned with the trace built before
dispatch. `Evaluate` is unchanged and fires no handlers (stays idempotent).

Each handler has signature `func(t Trace) error`. The single `Trace` argument is the same
read-only value `Run` returns, with exactly three methods:

| Method | File | Returns |
|---|---|---|
| `t.All()` | trace_all.go | `[]TraceEntry` — every node's entry in registration order |
| `t.Get(name)` | trace_get.go | `(TraceEntry, bool)` — one node's entry by short name |
| `t.Ctx()` | trace_ctx.go | `Context` — this Run's context |

`self` is not a separate argument: a handler attached with `g.Rule(X).RunOn(h)` already knows
its own node is `X`, so it finds its own `TraceEntry` via `t.Get("X")`. `t.All()` is the whole
graph's final state, useful for audit logging, explanations, and gradient thresholds via
`t.All()[i].Verdict`; `t.Ctx()` exposes this Run's `Context` typed (each `TraceEntry.Ground`
also holds it, but as `any`). To inspect another node's outcome, filter `t.All()`:
`Verdict > 0` prevailed, `Verdict <= 0` defeated, `Activated == false` did not apply. There is
no direct attacker lookup — `TraceEntry` carries no defeat-edge info, so the former
`RunView.Attackers(name)` is gone; reason from `t.All()[i].Verdict`. `ctx` is mutable; the
trace is the judgment record. Since `TraceEntry.Ground = ctx`, keep `ctx` to serialisable
values if you `json.Marshal` the trace.

#### Execution composition

`rule.Run(g)` declares an execution edge: when the node is **Active**, its sub-graph `g` is
Run with the same `ctx` — a graph-of-graphs. It is the execution counterpart of `Attacks`
(a defeat edge): `Attacks` composes *judgment* (verdict flows up into the target), `Run`
composes *execution* (an Active node drives a child graph down). For an Active node the
`RunOn` handler fires first, then the sub-graph is Run; `RunOn` and `Run(g)` coexist.
The sub-graph's verdicts stay isolated — only errors propagate (wrapped `run "node" → "subgraph": ...`)
— and each level builds its own trace. Execution composition must be a DAG: `Run` rejects
a static cycle once at the top-level entry via `detectRunCycle` (a shared sub-graph reached by
two paths is legal), and a runtime depth guard (`runMaxDepth = 64`) backstops runaway recursion.
`Run(nil)` panics.

### Testing Helper

| Function | File | Description |
|---|---|---|
| `RunCases(t, g, cases)` | run_cases.go | Run table-driven test cases against a graph |
| `VerdictAbove(v)` | verdict_above.go | Expectation: verdict > v |
| `VerdictAtMost(v)` | verdict_at_most.go | Expectation: verdict <= v |
| `VerdictBetween(lo, hi)` | verdict_between.go | Expectation: lo < verdict <= hi |
| `NoResult` | no_result.go | Expectation: no active warrants |

### Utilities

| Function | File | Description |
|---|---|---|
| `DetectCycle(edges)` | detect_cycle.go | DFS cycle detection |
| `FuncName(fn)` | func_name.go | Extract short function name |

## Types

| Type | File | Description |
|---|---|---|
| `Spec` | spec.go | Interface: `SpecName() string`, `Validate() error` |
| `Specs` | specs.go | Collection of Spec (`[]Spec`) with `Find` lookup |
| `EvalOption` | eval_option.go | Evaluation options (Method, Trace, Duration) |
| `EvalMethod` | eval_method.go | Verdict computation method (Matrix) |
| `Graph` | graph.go | Graph builder |
| `Rule` | rule.go | Opaque rule reference |
| `RuleMeta` | rule_meta.go | Rule metadata (Name, Qualifier, Strength, Defeats, Specs, Fn, RunOn, RunGraph) |
| `NodeHandler` | node_handler.go | Handler signature `func(t Trace) error` |
| `Trace` | trace.go | Read-only view of one Run: all node entries + ctx (`All()`, `Get(name)`, `Ctx()`) |
| `EvalResult` | eval_result.go | Verdict + trace |
| `TraceEntry` | trace_entry.go | Single rule evaluation record (Name=Claim, Ground=ctx, Specs=Backing, Verdict) |
| `Strength` | strength.go | Defeasible / Strict / Defeater |
| `TestCase` | test_case.go | Table-driven test case (Name, Context, Option, Expect) |
| `Context` | context.go | Interface: `Get(key string) (any, bool)`, `Set(key string, value any)` |
| `MapContext` | map_context.go | Default Context implementation |
| `Expectation` | expectation.go | Verdict assertion function `func([]EvalResult) error` |

## h-Categoriser

```
raw(a) = w(a) / (1 + Σ raw(attackers))     [0, 1]
verdict(a) = 2 × raw(a) - 1                [-1, 1]
```

+1.0 confirmed, 0.0 undecided, -1.0 fully rebutted.

## Internal Files

| File | Description |
|---|---|
| validate_spec.go | Validate Spec interface implementation |
| validate_spec_fields.go | Reject func fields in Spec structs |
| eval_context.go | Shared state for lazy evaluation |
| eval_context_calc.go | h-Categoriser recursive verdict |
| eval_context_calc_trace.go | h-Categoriser with trace and optional duration |
| eval_context_eval_rule.go | Evaluate single rule function |
| eval_context_record_trace.go | Record trace entry with optional duration measurement |
| eval_context_reset.go | Reset state between warrants |
| eval_method.go | EvalMethod type (Matrix) |
| resolve_option.go | Resolve EvalOption from variadic args |
| new_eval_context.go | Build eval context from graph |
| build_attacker_set.go | Build set of all attacker nodes |
| build_edges_from_rules.go | Build edge map from RuleMeta |
| collect_trace.go | Collect trace entries for evaluation |
| defeat_edge.go | Internal defeat edge struct |
| func_id.go | Full path function identifier |
| rule_id.go | Rule identifier (funcID + spec) |
| short_name.go | Extract short name from full path |
| infer_role.go | Infer rule role from graph position |
| is_warrant.go | Check if node is a warrant |
| to_rule_func.go | Convert any to rule function |
| safe_call.go | Panic-safe rule function invocation |
| safe_call_handler.go | Panic-safe node handler invocation |
| graph_evaluate_internal.go | Shared evaluation for Evaluate (lazy) and Run (full pass) |
| eval_context_fill_all.go | Calc every not-yet-run node to fill inactive state (full pass) |
| graph_run_depth.go | Recursive Run dispatch: assemble the flat trace, fire each Active node's RunOn handler + Run its sub-graph, with depth guard (runMaxDepth) |
| detect_run_cycle.go | 3-color DFS over RunGraph edges rejecting cyclic execution composition (must be a DAG) |
| specs_find.go | Find a Spec by type from Specs collection |
