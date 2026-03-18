//ff:type feature=graph type=model
//ff:what RuleDef — single rule entry in YAML graph definition
package graphdef

// RuleDef represents a single rule entry in the YAML definition.
type RuleDef struct {
	Name      string  `yaml:"name"`
	Role      string  `yaml:"role"`
	Qualifier float64 `yaml:"qualifier"`
}
