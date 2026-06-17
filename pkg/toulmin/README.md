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
| `g.Run(ctx, opt ...EvalOption)` | graph_run.go | Pre-evaluate (full pass), fire each node's event handler with a shared `RunView`, and Run any Active node's sub-graph; returns `([]EvalResult, RunView, error)` |
| `rule.OnActive(h)` | rule_on_active.go | Register handler fired when node is Active (func true && verdict > 0) |
| `rule.OnDefeated(h)` | rule_on_defeated.go | Register handler fired when node is Defeated (func true && verdict <= 0) |
| `rule.OnInactive(h)` | rule_on_inactive.go | Register handler fired when node is Inactive (func false) |
| `rule.Run(g)` | rule_run.go | Declare execution edge: when this node is Active, Run sub-graph `g` (graph-of-graphs; panics on nil) |
| `NewContext()` | new_context.go | Create `*MapContext` implementing Context interface |

`Run` forces Trace/Duration off and does a full pass so every node fires exactly one of
Inactive / Active / Defeated (verdict == 0 counts as Defeated). Handlers fire in rule
registration order; a handler error or panic stops `Run` immediately and is returned with
the `RunView` snapshot built before dispatch. `Evaluate` is unchanged and fires no handlers
(stays idempotent).

Before any handler fires, `Run` builds one immutable `RunView` snapshot of every node's
final event and shares it with all handlers. A handler reads its own `ev` plus the whole
graph's final state via `view.All()` / `view.Get(name)` / `view.Attackers(name)` — useful
for audit logging, explanations, and gradient thresholds (`view.Get(...).Verdict`). The
snapshot is immutable: mutating `ctx` in one handler never changes the `view` another sees
(`Type`/`Verdict` are copied; `Evidence` shares a reference and must be treated read-only).

#### Execution composition

`rule.Run(g)` declares an execution edge: when the node is **Active**, its sub-graph `g` is
Run with the same `ctx` — a graph-of-graphs. It is the execution counterpart of `Attacks`
(a defeat edge): `Attacks` composes *judgment* (verdict flows up into the target), `Run`
composes *execution* (an Active node drives a child graph down). For an Active node the
`OnActive` handler fires first, then the sub-graph is Run; `OnActive` and `Run(g)` coexist.
The sub-graph's verdicts stay isolated — only errors propagate (wrapped `run "node" → "subgraph": ...`)
— and each level builds its own `RunView`. Execution composition must be a DAG: `Run` rejects
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
| `RuleMeta` | rule_meta.go | Rule metadata (Name, Qualifier, Strength, Defeats, Specs, Fn, OnActive/OnDefeated/OnInactive, RunGraph) |
| `NodeEventType` | node_event_type.go | Node event classification (Inactive / Active / Defeated) |
| `NodeEvent` | node_event.go | Handler payload (Name, Role, Type, Verdict, Evidence) |
| `NodeHandler` | node_handler.go | Handler signature `func(ctx Context, ev NodeEvent, view RunView) error` |
| `RunView` | run_view.go | Read-only snapshot of every node's final event: `All()`, `Get(name)`, `Attackers(name)` |
| `EvalResult` | eval_result.go | Verdict + trace |
| `TraceEntry` | trace_entry.go | Single rule evaluation record |
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
| classify_event.go | Classify a node event from active flag and verdict |
| select_handler.go | Pick the handler matching a node event type |
| graph_run_depth.go | Recursive Run dispatch: fire handlers + Run Active nodes' sub-graphs, with depth guard (runMaxDepth) |
| new_run_view.go | Build the immutable RunView snapshot from a full-pass eval context |
| detect_run_cycle.go | 3-color DFS over RunGraph edges rejecting cyclic execution composition (must be a DAG) |
| node_event_type_string.go | NodeEventType.String() for logs and errors |
| specs_find.go | Find a Spec by type from Specs collection |
