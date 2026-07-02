//ff:func feature=tangl type=codegen control=iteration dimension=1
//ff:what buildInternals — renders every tangl:Internal entry in document order
package gen

import (
	"strings"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// buildInternals renders every tangl:Internal entry — an "on" event
// handler or an "every" tick runner — in document order.
func buildInternals(w *strings.Builder, doc *ast.Document) {
	for i, in := range doc.Internals {
		if in.Kind == ast.EveryTick {
			renderInternalEvery(w, in, i)
		} else {
			renderInternalOn(w, in, i)
		}
	}
}
