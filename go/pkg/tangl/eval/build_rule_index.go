//ff:func feature=tangl type=engine control=iteration dimension=1
//ff:what buildRuleIndex — builds the per-document rule lookup tables
package eval

import "github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"

// buildRuleIndex builds the ruleIndex lookup tables from a document's rules in
// document order.
func buildRuleIndex(doc *ast.Document) ruleIndex {
	idx := ruleIndex{
		ruleSet: map[string]bool{},
		lineOf:  map[string]int{},
		byName:  map[string]ast.InlineRule{},
		names:   make([]string, 0, len(doc.Rules)),
	}
	for _, r := range doc.Rules {
		idx.ruleSet[r.Name] = true
		idx.lineOf[r.Name] = r.Line
		idx.byName[r.Name] = r
		idx.names = append(idx.names, r.Name)
	}
	return idx
}
