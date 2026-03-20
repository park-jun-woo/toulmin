//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what LoadGraph — builds a live Graph from a GraphDef and function/backing registries
package toulmin

import "fmt"

// LoadGraph builds a *Graph from a GraphDef, a function registry, and an optional backing registry.
// Functions maps rule names to rule functions. Backings maps rule names to backing values (nil if absent).
// Returns an error if a rule name is not found in the registry or a defeat edge references an unknown rule.
func LoadGraph(def GraphDef, functions map[string]any, backings map[string]any) (*Graph, error) {
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
		default:
			return nil, fmt.Errorf("toulmin: rule %q has unknown role %q", rd.Name, rd.Role)
		}
		refs[rd.Name] = rule
	}

	for _, ed := range def.Defeats {
		from, ok := refs[ed.From]
		if !ok {
			return nil, fmt.Errorf("toulmin: defeat edge from %q not found", ed.From)
		}
		to, ok := refs[ed.To]
		if !ok {
			return nil, fmt.Errorf("toulmin: defeat edge to %q not found", ed.To)
		}
		g.Defeat(from, to)
	}

	return g, nil
}
