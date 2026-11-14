//ff:func feature=tangl type=validator control=iteration dimension=2 pattern=error-collection
//ff:what checkWeights — verifies every node qualifier and explicit attack force is a weight in [0, 1] at 4 decimals
package validate

import "github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"

// checkWeights verifies every case node's qualifier (`qualified`/`가중치`) and
// every explicit attack force (`with force`/`세기`) with ast.CheckUnitInterval:
// NaN, a value outside [0, 1], or one finer than four decimal places is a load
// error. The surface parsers already reject such literals; this check covers an
// AST that reached the toolchain another way (a .geul stream, a hand-built AST).
func checkWeights(doc *ast.Document) []error {
	var errs []error
	for _, c := range doc.Cases {
		for _, n := range c.Nodes {
			if n.Qualified == nil {
				continue
			}
			if err := ast.CheckUnitInterval(*n.Qualified); err != nil {
				errs = append(errs, errAt(doc.Path, n.Line, "case %q: node %q qualified %v: %v", c.Name, n.Name, *n.Qualified, err))
			}
		}
		for _, a := range c.Attacks {
			if a.Force == nil {
				continue
			}
			if err := ast.CheckUnitInterval(*a.Force); err != nil {
				errs = append(errs, errAt(doc.Path, a.Line, "case %q: attack %q -> %q force %v: %v", c.Name, a.Attacker, a.Target, *a.Force, err))
			}
		}
	}
	return errs
}
