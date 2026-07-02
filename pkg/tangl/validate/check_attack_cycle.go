//ff:func feature=tangl type=validator control=iteration dimension=1
//ff:what checkAttackCycle — rejects a circular don't/do-not defeat graph within a case
package validate

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// checkAttackCycle rejects a cyclic `don't <target> when <attacker>` defeat
// graph within a single case (node A attacking node B attacking ... node A).
func checkAttackCycle(doc *ast.Document) []error {
	var errs []error
	for _, c := range doc.Cases {
		edges := caseAttackEdges(c)
		lines := caseNodeLineIndex(c)
		if err := detectNameCycle(doc.Path, "don't", edges, func(name string) int { return lines[name] }); err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}
