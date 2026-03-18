//ff:type feature=graph type=model
//ff:what GraphDef — YAML graph definition (graph name, rules, defeats)
package graphdef

// GraphDef represents a graph definition parsed from YAML.
type GraphDef struct {
	Graph   string    `yaml:"graph"`
	Rules   []RuleDef `yaml:"rules"`
	Defeats []EdgeDef `yaml:"defeats"`
}
