//ff:type feature=graph type=model
//ff:what RuleDef — single rule entry in YAML graph definition
package graphdef

// RuleDef represents a single rule entry in the YAML definition.
// Qualifier is a pointer to distinguish unset (nil → default 1.0) from explicit 0.0.
type RuleDef struct {
	Name      string   `yaml:"name"`
	Role      string   `yaml:"role"`
	Qualifier *float64 `yaml:"qualifier"`
}
