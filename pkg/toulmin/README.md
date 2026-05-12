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
| `NewContext()` | new_context.go | Create `*MapContext` implementing Context interface |

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
| `RuleMeta` | rule_meta.go | Rule metadata (Name, Qualifier, Strength, Defeats, Specs, Fn) |
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
| specs_find.go | Find a Spec by type from Specs collection |
