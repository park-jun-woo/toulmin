//ff:func feature=tangl type=codegen control=iteration dimension=2
//ff:what needsCaseHelper — reports whether the file needs the shared tanglCaseActive helper
package gen

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// needsCaseHelper reports whether the generated file needs the shared
// tanglCaseActive helper: any node "checking" another case, or any
// Internal "every ... until" clause, composes a case's Evaluate results
// into a single verdict this way.
func needsCaseHelper(doc *ast.Document) bool {
	for _, c := range doc.Cases {
		for _, n := range c.Nodes {
			if n.Checking != "" {
				return true
			}
		}
	}
	for _, in := range doc.Internals {
		if in.Until != "" {
			return true
		}
	}
	return false
}
