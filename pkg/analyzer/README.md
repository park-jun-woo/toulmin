# pkg/analyzer

Go AST analysis for toulmin defeats graph extraction. Parses Go source files and extracts `NewGraph` → `Warrant/Rebuttal/Defeater` → `Defeat` call chains without compilation.

## Public API

| Function | File | Description |
|---|---|---|
| `ExtractDefeats(path)` | extract_defeats.go | Go source file → `[]DefeatGraph` |

## Types

| Type | File | Description |
|---|---|---|
| `DefeatGraph` | defeat_graph.go | Extracted graph (Name, Rules, Defeats map) |

## Internal Files

| File | Description |
|---|---|
| find_graph_calls.go | Find `NewGraph` CallExpr in var declarations |
| extract_graph_name.go | Extract graph name from `NewGraph("name")` |
| collect_chain.go | Walk method chains collecting rule registrations and defeat edges |
| collect_rule_name.go | Extract rule name from `Warrant/Rebuttal/Defeater` call |
| collect_defeat_edge.go | Extract defeat edge from `Defeat(from, to)` call |
| collect_value_specs.go | Collect AST ValueSpec nodes |

## Usage

```go
graphs, err := analyzer.ExtractDefeats("policy.go")
for _, g := range graphs {
    // g.Name     = graph name
    // g.Rules    = []string of rule names
    // g.Defeats  = map[target][]attacker
    err := toulmin.DetectCycle(g.Defeats)
}
```

## Flow

```
Go source file → AST parse → find NewGraph calls → walk method chains → DefeatGraph
```
