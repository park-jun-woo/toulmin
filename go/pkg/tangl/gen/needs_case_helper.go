//ff:func feature=tangl type=codegen control=iteration dimension=1
//ff:what needsCaseHelper — reports whether the file needs the shared tanglCaseActive helper
package gen

import "github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"

// needsCaseHelper reports whether the generated file needs the shared
// tanglCaseActive helper. Every case now emits an exported checking wrapper
// (buildCheckingWrappers) that composes that case's Evaluate results through
// this helper, so any document with at least one case needs it; an Internal
// "every ... until" clause also composes results this way.
func needsCaseHelper(doc *ast.Document) bool {
	if len(doc.Cases) > 0 {
		return true
	}
	for _, in := range doc.Internals {
		if in.Until != "" {
			return true
		}
	}
	return false
}
