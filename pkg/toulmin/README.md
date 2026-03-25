# pkg/toulmin

Toulmin argumentation-based rule engine core. Defeats graph + h-Categoriser verdict.

## Rule Function Signature

```go
func(ctx Context, backing Backing) (bool, any)
// ctx = Context interface with Get/Set for per-request facts
// backing = fixed criteria at declaration
// returns (judgment, evidence)
```

## Public API

### Graph Builder

| Function | File | Description |
|---|---|---|
| `NewGraph(name)` | new_graph.go | Create graph builder |
| `g.Rule(fn)` | graph_rule.go | Register rule, return `*Rule` (builder: `.Backing(b)`, `.Qualifier(q)`) |
| `g.Counter(fn)` | graph_counter.go | Register counter, return `*Rule` (builder: `.Backing(b)`, `.Qualifier(q)`) |
| `g.Except(fn)` | graph_except.go | Register except, return `*Rule` (builder: `.Backing(b)`, `.Qualifier(q)`) |
| `rule.Attacks(target)` | rule_attacks.go | Declare defeat edge (method on `*Rule`) |
| `g.Evaluate(ctx, opt ...EvalOption)` | graph_evaluate.go | Run evaluation, return verdicts (ctx is Context, opt controls Trace/Duration) |
| `NewContext()` | new_context.go | Create `*MapContext` implementing Context interface |

### Dynamic Loading

| Function | File | Description |
|---|---|---|
| `ParseYAML(path)` | parse_yaml.go | YAML file → `GraphDef` |
| `ValidateGraphDef(def)` | validate_graph_def.go | Check graph name, edge refs, cycles |
| `LoadGraph(def, funcs, backings map[string]Backing)` | load_graph.go | `GraphDef` → live `*Graph` |

### Engine (legacy)

| Function | File | Description |
|---|---|---|
| `NewEngine()` | new_engine.go | Create legacy engine |
| `e.Register(RuleMeta)` | engine_register.go | Register rule by name |
| `e.Evaluate(ctx, opt ...EvalOption)` | engine_evaluate.go | Run evaluation (ctx is Context, opt controls Trace/Duration) |

### Code Generation

| Function | File | Description |
|---|---|---|
| `GenerateGraph(pkg, def)` | generate_graph.go | `GraphDef` → Go source code |

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
| `Backing` | backing.go | Interface: `BackingName() string`, `Validate() error` |
| `EvalOption` | eval_option.go | Evaluation options (Method, Trace, Duration) |
| `EvalMethod` | eval_method.go | Verdict computation method (Matrix, Recursive) |
| `Graph` | graph.go | Graph builder |
| `Rule` | rule.go | Opaque rule reference |
| `RuleMeta` | rule_meta.go | Rule metadata (Name, Qualifier, Strength, Defeats, Backing, Fn) |
| `EvalResult` | eval_result.go | Verdict + trace |
| `TraceEntry` | trace_entry.go | Single rule evaluation record |
| `Strength` | strength.go | Defeasible / Strict / Defeater |
| `Engine` | engine.go | Legacy rule registry |
| `GraphDef` | graph_def.go | Graph definition (AST from YAML/DB/API) |
| `GraphRuleDef` | graph_rule_def.go | Rule entry in GraphDef |
| `GraphEdgeDef` | graph_edge_def.go | Defeat edge in GraphDef |
| `TestCase` | test_case.go | Table-driven test case (Name, Ctx, Option, Expect) |
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
| resolve_backing.go | Resolve backing from registry for LoadGraph |
| validate_backing.go | Validate Backing interface implementation |
| validate_backing_fields.go | Reject func fields in Backing structs |
| eval_context.go | Shared state for lazy evaluation |
| eval_context_calc.go | h-Categoriser recursive verdict |
| eval_context_calc_trace.go | h-Categoriser with trace and optional duration |
| eval_context_record_trace.go | Record trace entry with optional duration measurement |
| eval_method.go | EvalMethod type (Matrix, Recursive) |
| resolve_option.go | Resolve EvalOption from variadic args |
| eval_context_reset.go | Reset state between warrants |
| new_eval_context.go | Build eval context from graph |
| build_attacker_set.go | Build set of all attacker nodes |
| build_edges_from_rules.go | Build edge map from RuleMeta |
| defeat_edge.go | Internal defeat edge struct |
| func_id.go | Full path function identifier |
| rule_id.go | Rule identifier (funcID + backing) |
| short_name.go | Extract short name from full path |
| infer_role.go | Infer rule role from graph position |
| is_warrant.go | Check if node is a warrant |
| to_rule_func.go | Convert any to rule function |
| wrap_legacy.go | Wrap legacy 2-arg signature |
| graph_var_name.go | Graph name → PascalCase variable name |
| role_to_method.go | YAML role → Graph method name |
| rule_var_name.go | Rule name → camelCase variable name |

## YAML Schema

```yaml
graph: <name>
rules:
  - name: <rule_name>
    role: rule | counter | except
    qualifier: 0.0-1.0  # default 1.0
defeats:
  - from: <attacker>
    to: <target>
```

## Flow

```
YAML file → ParseYAML() → GraphDef → ValidateGraphDef() → LoadGraph(def, funcs, backings) → *Graph → Evaluate()
                                    → GenerateGraph(pkg, def) → Go source code
```
