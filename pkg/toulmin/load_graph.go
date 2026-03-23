//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what LoadGraph — validates GraphDef and builds a live Graph
package toulmin

import "fmt"

// LoadGraph validates a GraphDef, then builds a *Graph using the provided function and backing registries.
// Returns an error if validation fails or a rule name is not found in the function registry.
func LoadGraph(def GraphDef, functions map[string]any, backings map[string]any) (*Graph, error) {
	if err := ValidateGraphDef(def); err != nil {
		return nil, err
	}
	g := NewGraph(def.Graph)
	refs := make(map[string]*Rule, len(def.Rules))

	for _, rd := range def.Rules {
		fn, ok := functions[rd.Name]
		if !ok {
			return nil, fmt.Errorf("toulmin: rule %q not found in function registry", rd.Name)
		}

		q := rd.Qualifier
		if q == 0 {
			q = 1.0
		}

		var backing any
		if backings != nil {
			backing = backings[rd.Name]
		}

		var rule *Rule
		switch rd.Role {
		case "warrant":
			rule = g.Warrant(fn, backing, q)
		case "rebuttal":
			rule = g.Rebuttal(fn, backing, q)
		case "defeater":
			rule = g.Defeater(fn, backing, q)
		}
		refs[rd.Name] = rule
	}

	for _, ed := range def.Defeats {
		g.Defeat(refs[ed.From], refs[ed.To])
	}

	return g, nil
}
