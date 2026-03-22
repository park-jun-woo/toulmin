//ff:type feature=engine type=model
//ff:what GraphRuleDef — single rule entry in graph definition
package toulmin

// GraphRuleDef defines a single rule entry.
type GraphRuleDef struct {
	Name      string  `yaml:"name"`
	Role      string  `yaml:"role"`      // "warrant", "rebuttal", "defeater"
	Qualifier float64 `yaml:"qualifier"` // 0 means default 1.0
}
