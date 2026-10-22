//ff:func feature=tangl type=validator control=iteration dimension=3
//ff:what checkWithRefs — verifies with-clause terms name an existing Definitions entry
package validate

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// checkWithRefs verifies that every `with <term>` argument on a node names a
// term declared in the document's tangl:Definitions section.
func checkWithRefs(doc *ast.Document) []error {
	defs := defNameSet(doc)
	defByName := make(map[string]ast.Definition, len(doc.Defs))
	for _, d := range doc.Defs {
		defByName[d.Name] = d
	}
	var errs []error
	for _, rule := range doc.Rules {
		explicit := make(map[string]bool, len(rule.With))
		for _, term := range rule.With {
			explicit[term] = true
		}
		terms := append(append([]string(nil), rule.With...), ast.RuleSpecNames(rule)...)
		seen := map[string]bool{}
		for _, term := range terms {
			if seen[term] {
				continue
			}
			seen[term] = true
			if !defs[term] && explicit[term] {
				errs = append(errs, errAt(doc.Path, rule.Line, "rule %q references undefined definition %q", rule.Name, term))
			} else if defs[term] && !specCapable(defByName[term]) {
				errs = append(errs, errAt(doc.Path, rule.Line, "rule %q references definition %q that cannot produce a Spec", rule.Name, term))
			}
		}
	}
	for _, c := range doc.Cases {
		for _, n := range c.Nodes {
			for _, term := range n.With {
				if !defs[term] {
					errs = append(errs, errAt(doc.Path, n.Line, "case %q: node %q references undefined term %q", c.Name, n.Name, term))
				}
			}
		}
	}
	return errs
}
