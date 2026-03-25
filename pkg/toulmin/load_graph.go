//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what LoadGraph — validates GraphDef and builds a live Graph
package toulmin

import "fmt"

// LoadGraph validates a GraphDef, then builds a *Graph using the provided function and backing registries.
// Returns an error if validation fails or a rule name is not found in the function registry.
func LoadGraph(def GraphDef, functions map[string]any, backings map[string]Backing) (*Graph, error) {
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

		backing, err := resolveBacking(rd.Name, backings)
		if err != nil {
			return nil, err
		}

		var rule *Rule
		switch rd.Role {
		case "rule":
			rule = g.Rule(fn)
		case "counter":
			rule = g.Counter(fn)
		case "except":
			rule = g.Except(fn)
		}
		if backing != nil {
			rule.Backing(backing)
		}
		if q != 1.0 {
			rule.Qualifier(q)
		}
		refs[rd.Name] = rule
	}

	for _, ed := range def.Defeats {
		refs[ed.From].Attacks(refs[ed.To])
	}

	return g, nil
}
