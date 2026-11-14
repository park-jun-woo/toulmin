//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what internalBody — find the Internal trigger declared at a source line and return its body slices
package en

import "github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"

// internalBody returns the Runs/Checks slices and index of the Internal trigger
// whose declaration is on source line no, or ok=false if none matches.
func internalBody(a *ast.Document, no int) (runs, checks []string, idx int, ok bool) {
	for i := range a.Internals {
		if a.Internals[i].Line == no {
			return a.Internals[i].Runs, a.Internals[i].Checks, i, true
		}
	}
	return nil, nil, 0, false
}
