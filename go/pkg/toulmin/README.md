# pkg/toulmin

Toulmin argumentation-based rule engine core. Defeats graph + h-Categoriser verdict.

## Rule Function Signature

```go
func(ctx Context, specs Specs) (float64, any)
// ctx = Context interface with Get/Set for per-request facts
// specs = fixed criteria attached via .With() at declaration
// returns (v, evidence), v in [-1, +1]: +1 fully applies, 0 half applies, -1 does not apply

func(ctx Context, specs Specs) (bool, any) // legacy, still accepted: true → +1, false → -1
```

**"Does not apply" is −1, not 0. Zero means "half applies".** If you carry bool habits over and
write `return 0, nil` for "no", the rule stays Active (applicability 0.5), its attackers are
evaluated, and it keeps weighing into the graph — the verdict silently comes out wrong. **The
engine cannot detect this, because 0 is a legal value; you must watch for it yourself.**

- NaN or a value outside `[-1, +1]` is a runtime error: `Evaluate`/`Run` return
  `rule "<id>": rule verdict <v> out of range [-1, 1]` (same path as a panicking rule; checked in `safeCall`).
- The bool signature is wrapped at registration (`toRuleFunc` → `boolRuleFunc`); the rule ID is
  still taken from the original function. Bool graphs give exactly the same verdicts as before,
  except where an `Except` attacks (see h-Categoriser).
- A node is Active when `v > -1`; only `v == -1` skips its attackers (lazy evaluation).
  `TraceEntry.RuleVerdict` records `v`, `TraceEntry.Activated` is `v > -1`.

## Public API

### Graph Builder

| Function | File | Description |
|---|---|---|
| `NewGraph(name)` | new_graph.go | Create graph builder |
| `g.Rule(fn)` | graph_rule.go | Register rule, return `*Rule` (builder: `.With(spec)`, `.Qualifier(q)`) |
| `g.Counter(fn)` | graph_counter.go | Register counter, return `*Rule` (builder: `.With(spec)`, `.Qualifier(q)`) |
| `g.Except(fn)` | graph_except.go | Register except, return `*Rule` (builder: `.With(spec)`, `.Qualifier(q)`) |
| `rule.Attacks(target, force ...float64)` | rule_attacks.go | Declare defeat edge (method on `*Rule`). Optional force `F` in `[0, 1]`, default `1`; more than one value, NaN, or out of range panics (like `Qualifier`) |
| `rule.Supersedes(target)` | rule_supersedes.go | Declare supersede edge: when this rule's applicability `a_eff > 0` (Except gates count, Counter rebuttals do not), target is crisply excluded from the tally (sets `TraceEntry.Excluded`); never enters the h-Categoriser traversal |
| `rule.With(spec)` | rule_with.go | Attach a Spec to the rule (additive, supports chaining) |
| `rule.Level(l)` | rule_level.go | Set the crisp reporting `Level` (`Unset`/`Fail`/`Review`); the explicit value (including `Unset`) overrides the role-based default (method on `*Rule`, chains) |
| `g.Evaluate(ctx, opt ...EvalOption)` | graph_evaluate.go | Run evaluation, return verdicts (ctx is Context, opt controls Trace/Duration) |
| `g.Run(ctx, opt ...EvalOption)` | graph_run.go | Pre-evaluate (full pass), fire each Active node's `RunOn` handler with the whole-graph `Trace`, and Run any Active node's sub-graph; returns `([]EvalResult, Trace, error)` |
| `rule.RunOn(h)` | rule_run_on.go | Register the handler fired when this node is Active (rule `v > -1` && verdict > 0) |
| `rule.Run(g)` | rule_run.go | Declare execution edge: when this node is Active, Run sub-graph `g` (graph-of-graphs; panics on nil) |
| `NewContext()` | new_context.go | Create `*MapContext` implementing Context interface |

`Run` forces Trace/Duration off and does a full pass, then fires the `RunOn` handler of every
**Active** node — a node that applied and prevailed (`Activated && Verdict > 0`). That is the
only event; there are no Defeated/Inactive handlers. Handlers fire in rule registration order;
a handler error or panic stops `Run` immediately and is returned with the trace built before
dispatch. `Evaluate` is unchanged and fires no handlers (stays idempotent).

Each handler has signature `func(self TraceEntry, t Trace) error`. `self` is the handler's own
node entry; the `Trace` argument `t` is the same read-only value `Run` returns, with exactly
three methods:

| Method | File | Returns |
|---|---|---|
| `t.All()` | trace_all.go | `[]TraceEntry` — every node's entry in registration order |
| `t.Get(name)` | trace_get.go | `(TraceEntry, bool)` — one node's entry by short name |
| `t.Ctx()` | trace_ctx.go | `Context` — this Run's context |

`self` is the first argument, handed in by the engine, so a handler never looks up its own row
by name (which was fragile when an entry's name carried a dynamic suffix). `t.All()` is the whole
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
| `NodeHandler` | node_handler.go | Handler signature `func(self TraceEntry, t Trace) error` |
| `Trace` | trace.go | Read-only view of one Run: all node entries + ctx (`All()`, `Get(name)`, `Ctx()`) |
| `EvalResult` | eval_result.go | Verdict + evidence + `Applicability` (`a_eff`, see h-Categoriser) + trace |
| `TraceEntry` | trace_entry.go | Single rule evaluation record (Name=Claim, Ground=ctx, Specs=Backing, Verdict; `RuleVerdict` = the rule's own `v`, `Activated` = `v > -1`; `Excluded` = crisply superseded out of the tally; `Level` = crisp reporting level) |
| `Strength` | strength.go | Defeasible / Strict / Defeater |
| `Level` | level.go | Crisp reporting level: `Unset` / `Fail` / `Review` (see below) |
| `Fact` | fact.go | Optional structured evidence convention `{Where, Expected, Actual any}` |
| `TestCase` | test_case.go | Table-driven test case (Name, Context, Option, Expect) |
| `Context` | context.go | Interface: `Get(key string) (any, bool)`, `Set(key string, value any)` |
| `MapContext` | map_context.go | Default Context implementation |
| `Expectation` | expectation.go | Verdict assertion function `func([]EvalResult) error` |

## h-Categoriser

```
a       = (v + 1) / 2                              own applicability, v = rule function's return
S       = Σ_counter F · (v_j + 1) / 2             rebuttals (Counter / Rule attackers): denominator
a_eff   = a · Π_except (1 − F · (v_j + 1) / 2)    exceptions (Except): multiplicative gate
raw(a)  = q · a_eff / (1 + S)                      [0, 1]
verdict = 2 × raw − 1                              [-1, 1]
```

`v_j` is the attacker's verdict and `F` the edge force (`Attacks(target, F)`, default 1).
+1.0 confirmed, 0.0 undecided, -1.0 fully rebutted.

- A Counter argues against the conclusion and adds to the denominator (unchanged).
- An Except says the rule does not apply here, so it multiplies the target's applicability.
  One fully applying exception gates the target to −1; three half-applying ones leave
  `(1 − 0.75)³ ≈ 0.016`. (Previously exceptions sat in the denominator, so three of them
  only pulled the target to −0.5.)
- Whether an attacker is an exception is read from its node (`Strength == Defeater` or role `except`).
- `Strict` nodes ignore both rebuttals and exceptions.
- `Supersedes` excludes its target only while the superseder's `a_eff > 0`: an exception that fully gates the superseder lifts the exclusion; a rebuttal does not.
- With bool rules and no Except attackers, `a = a_eff = 1`, `F = 1`, and the result is
  bit-for-bit the former `raw = q / (1 + Σ raw(attackers))`.
- `EvalResult.Applicability` reports `a_eff` for each warrant.

## Level — crisp reporting

`Level` (`level.go`) is a crisp per-rule label orthogonal to the `[-1, 1]` verdict, for
downstream aggregation such as reins (any `Fail` → FAIL, else any `Review` → REVIEW, else PASS).
Values: `Unset` (zero value, "no level to report"), `Fail`, `Review`.

Each rule's **effective Level** — the value written to `TraceEntry.Level` on both the `Run` and
`Evaluate`-with-`Trace` paths — is resolved as:

- explicit `rule.Level(l)` was called → that value wins (even `Level(Unset)`, which opts out of the default);
- otherwise → the role-based default (`default_level_for_role.go`): `counter` → `Fail`,
  `except` → `Review`, everything else (including `rule`) → `Unset`.

`TraceEntry.Level` is serialised **without** `omitempty`: `Unset` (0) is a meaningful "no level"
report value, not an absent field.

## Fact — structured evidence

`Fact{Where, Expected, Actual any}` (`fact.go`) is an optional convention for the `any` evidence
a rule function returns as its second value. The engine never requires or interprets it — it flows
through unchanged into `TraceEntry.Evidence` / `EvalResult.Evidence`, and consumers type-assert
`evidence.(Fact)` to recover a located, quantified violation. Recovery per rule ID is already
provided by `Trace.All()` / `Trace.Get(name)` (whole graph) and `EvalResult.Name` + `.Evidence`
(per warrant) — no new API is needed. Determinism of `(verdict, root-cause rule ID, Fact)` holds
as long as each rule function is itself pure/deterministic (a documented precondition; not enforced).

## Internal Files

| File | Description |
|---|---|
| validate_spec.go | Validate Spec interface implementation |
| validate_spec_fields.go | Reject func fields in Spec structs |
| eval_context.go | Shared state for lazy evaluation |
| eval_context_calc.go | h-Categoriser recursive verdict |
| eval_context_combine.go | Fold attackers into a node's verdict (rebuttal denominator, exception gate, force) |
| eval_context_settle.go | `raw = q·a_eff/(1+S)`, record applicability, cache verdict |
| eval_context_set_outcome.go | Store a rule's own `v`, derived active flag, evidence |
| eval_context_edge_force.go | Force of an attack edge (default 1) |
| eval_context_is_except_attacker.go | Whether an attacker gates as an exception |
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
| default_level_for_role.go | Default crisp Level per role (counter→Fail, except→Review, else Unset) |
| eval_context_effective_level.go | Resolve a node's effective Level (explicit set wins over role default) |
| is_warrant.go | Check if node is a warrant |
| to_rule_func.go | Convert any to rule function (float64 as-is, bool wrapped) |
| bool_rule_func.go | Wrap a bool rule function (true → +1, false → -1) |
| attack_force.go | Resolve and validate `Attacks`' optional force |
| safe_call.go | Panic-safe rule function invocation with verdict range check |
| safe_call_handler.go | Panic-safe node handler invocation |
| graph_evaluate_internal.go | Shared evaluation for Evaluate (lazy) and Run (full pass) |
| eval_context_fill_all.go | Calc every not-yet-run node to fill inactive state (full pass) |
| graph_run_depth.go | Recursive Run dispatch: assemble the flat trace, fire each Active node's RunOn handler + Run its sub-graph, with depth guard (runMaxDepth) |
| detect_run_cycle.go | 3-color DFS over RunGraph edges rejecting cyclic execution composition (must be a DAG) |
| specs_find.go | Find a Spec by type from Specs collection |
