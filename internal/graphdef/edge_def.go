//ff:type feature=graph type=model
//ff:what EdgeDef — single defeat edge in YAML graph definition
package graphdef

// EdgeDef represents a single defeat edge in the YAML definition.
type EdgeDef struct {
	From string `yaml:"from"`
	To   string `yaml:"to"`
}
