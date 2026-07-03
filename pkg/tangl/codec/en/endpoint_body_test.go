//ff:func feature=tangl type=parser control=sequence
//ff:what TestEndpointBody — find the Provides endpoint declared at a source line, or report no match
package en

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// TestEndpointBody covers both outcomes: a line matching one of a.Provides
// returns its Runs/Checks and index, and a line matching none returns ok=false.
func TestEndpointBody(t *testing.T) {
	a := &ast.Document{Provides: []ast.Endpoint{
		{Name: "a", Line: 2, Runs: []string{"x"}},
		{Name: "b", Line: 5, Checks: []string{"y"}},
	}}

	runs, checks, idx, ok := endpointBody(a, 5)
	if !ok || idx != 1 || len(checks) != 1 || checks[0] != "y" || len(runs) != 0 {
		t.Errorf("endpointBody(a, 5) = (%v, %v, %d, %v), want ([], [y], 1, true)", runs, checks, idx, ok)
	}

	if _, _, _, ok := endpointBody(a, 99); ok {
		t.Error("endpointBody matched a line with no Provides declaration")
	}
}
