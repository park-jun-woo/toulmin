//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what endpointBody — find the Provides endpoint declared at a source line and return its body slices
package en

import "github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"

// endpointBody returns the Runs/Checks slices and index of the Provides endpoint
// whose declaration is on source line no, or ok=false if none matches.
func endpointBody(a *ast.Document, no int) (runs, checks []string, idx int, ok bool) {
	for i := range a.Provides {
		if a.Provides[i].Line == no {
			return a.Provides[i].Runs, a.Provides[i].Checks, i, true
		}
	}
	return nil, nil, 0, false
}
