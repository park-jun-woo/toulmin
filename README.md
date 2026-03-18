# toulmin

Toulmin argumentation-based rule engine for Go.

Rules are Go functions. The engine builds a defeats graph from activated rules and computes verdicts via Amgoud's h-Categoriser on a [-1, +1] scale.

## Install

```bash
go get github.com/park-jun-woo/toulmin/pkg/toulmin
```

## Usage

```go
eng := toulmin.NewEngine()

eng.Register(toulmin.RuleMeta{
    Name:      "OneFileOneFunc",
    Qualifier: 1.0,
    Strength:  toulmin.Defeasible,
    Backing:   "Bohm-Jacopini theorem",
    What:      "F1: one func per file",
    Fn:        func(claim any, ground any) bool { return ground.(FileGround).FuncCount > 1 },
})

eng.Register(toulmin.RuleMeta{
    Name:     "TestFileException",
    Qualifier: 1.0,
    Strength:  toulmin.Defeater,
    Defeats:   []string{"OneFileOneFunc"},
    Backing:   "test files conventionally group multiple test funcs",
    What:      "F5: test files allow multiple funcs",
    Fn:        func(claim any, ground any) bool { return strings.HasSuffix(claim.(string), "_test.go") },
})

results := eng.Evaluate(claim, ground)
// results[0].Verdict: +1.0 violation, 0.0 undecided, -1.0 rebutted
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

## Annotations

Rules declare metadata via `//rule:` comments:

```go
//rule:warrant qualifier=1.0 strength=strict
//rule:backing "Bohm-Jacopini theorem"
//rule:what F1: one func per file
func CheckOneFileOneFunc(claim any, ground any) bool { ... }

//rule:defeater defeats=CheckOneFileOneFunc
//rule:backing "test files conventionally group multiple test funcs"
//rule:what F5: test files allow multiple funcs
func TestFileException(claim any, ground any) bool { ... }
```

## Theory

| Component | Source |
|---|---|
| 6-element structure | Toulmin (1958) |
| strict/defeasible/defeater | Nute (1994) |
| h-Categoriser | Amgoud & Ben-Naim (2013, 2017) |

## License

MIT
