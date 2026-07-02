//ff:func feature=tangl type=codegen control=iteration dimension=1
//ff:what needsTime — reports whether the document has any Internal "every" trigger
package gen

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// needsTime reports whether the document has any Internal "every" tick
// trigger, which codegens a time.Ticker-driven runner.
func needsTime(doc *ast.Document) bool {
	for _, in := range doc.Internals {
		if in.Kind == ast.EveryTick {
			return true
		}
	}
	return false
}
