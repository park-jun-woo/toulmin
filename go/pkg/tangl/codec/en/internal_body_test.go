//ff:func feature=tangl type=parser control=sequence
//ff:what TestInternalBody — find the Internal trigger declared at a source line, or report no match
package en

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// TestInternalBody covers both outcomes: a line matching one of a.Internals
// returns its Runs/Checks and index, and a line matching none returns ok=false.
func TestInternalBody(t *testing.T) {
	a := &ast.Document{Internals: []ast.Internal{
		{Line: 3, Runs: []string{"x"}},
		{Line: 8, Checks: []string{"y"}},
	}}

	runs, checks, idx, ok := internalBody(a, 8)
	if !ok || idx != 1 || len(checks) != 1 || checks[0] != "y" || len(runs) != 0 {
		t.Errorf("internalBody(a, 8) = (%v, %v, %d, %v), want ([], [y], 1, true)", runs, checks, idx, ok)
	}

	if _, _, _, ok := internalBody(a, 99); ok {
		t.Error("internalBody matched a line with no Internal declaration")
	}
}
