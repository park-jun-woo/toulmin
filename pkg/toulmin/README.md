# pkg/toulmin

Toulmin argumentation-based rule engine core. Defeats graph + h-Categoriser verdict.

## Rule Function Signature

```go
func(claim any, ground any, backing any) (bool, any)
// ground = per-request facts, backing = fixed criteria at declaration
// returns (judgment, evidence)
```

## Public API

### Graph Builder

| Function | File | Description |
|---|---|---|
| `NewGraph(name)` | new_graph.go | Create graph builder |
| `g.Warrant(fn, backing Backing, qualifier)` | graph_warrant.go | Register warrant, return `*Rule` |
| `g.Rebuttal(fn, backing Backing, qualifier)` | graph_rebuttal.go | Register rebuttal, return `*Rule` |
| `g.Defeater(fn, backing Backing, qualifier)` | graph_defeater.go | Register defeater, return `*Rule` |
| `g.Defeat(from, to)` | graph_defeat.go | Declare defeat edge |
| `g.Evaluate(claim, ground, opt ...EvalOption)` | graph_evaluate.go | Run evaluation, return verdicts (opt controls Trace/Duration) |

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
| `e.Evaluate(claim, ground, opt ...EvalOption)` | engine_evaluate.go | Run evaluation (opt controls Trace/Duration) |

### Code Generation

| Function | File | Description |
|---|---|---|
| `GenerateGraph(pkg, def)` | generate_graph.go | `GraphDef` → Go source code |

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
    role: warrant | rebuttal | defeater
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
