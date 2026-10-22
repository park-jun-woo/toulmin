//ff:func feature=tangl type=parser control=sequence
//ff:what bodySlices — dispatch a body-line lookup to the endpoint or internal AST slices
package en

import (
	"github.com/geul-org/geul/go/pkg/toulmin"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// bodySlices returns the Runs/Checks slices and declaration index of the Provides
// endpoint or Internal trigger declared at source line no, dispatching on the
// section code. ok is false when no declaration matches that line.
func bodySlices(a *ast.Document, code toulmin.SectionCode, no int) (runs, checks []string, idx int, ok bool) {
	if code == toulmin.SectionProvides {
		return endpointBody(a, no)
	}
	return internalBody(a, no)
}
