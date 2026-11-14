//ff:func feature=tangl type=parser control=sequence
//ff:what TestBodySlices — dispatch a body lookup to endpointBody for Provides, internalBody otherwise
package en

import (
	"testing"

	"github.com/geul-org/gogeul/pkg/toulmin"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// TestBodySlices checks that SectionProvides dispatches to endpointBody (matching
// a.Provides by line) and any other section code (here SectionInternal) dispatches
// to internalBody (matching a.Internals by line).
func TestBodySlices(t *testing.T) {
	a := &ast.Document{
		Provides:  []ast.Endpoint{{Name: "ep", Line: 2, Runs: []string{"a"}}},
		Internals: []ast.Internal{{Line: 5, Runs: []string{"b"}}},
	}

	runs, _, idx, ok := bodySlices(a, toulmin.SectionProvides, 2)
	if !ok || idx != 0 || len(runs) != 1 || runs[0] != "a" {
		t.Errorf("Provides dispatch: (%v, %d, %v), want (true, 0, [a])", ok, idx, runs)
	}

	runs, _, idx, ok = bodySlices(a, toulmin.SectionInternal, 5)
	if !ok || idx != 0 || len(runs) != 1 || runs[0] != "b" {
		t.Errorf("Internal dispatch: (%v, %d, %v), want (true, 0, [b])", ok, idx, runs)
	}
}
