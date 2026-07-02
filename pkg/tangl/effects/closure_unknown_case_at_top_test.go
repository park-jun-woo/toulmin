//ff:func feature=tangl type=analyzer control=sequence
//ff:what closureUnknownCaseAtTop — verifies Closure errors when an endpoint's Runs list names a case absent from doc.Cases
package effects

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// closureUnknownCaseAtTop verifies that an endpoint whose Runs list names a
// case absent from doc.Cases surfaces the findCase-nil error from the top
// loop (visit called directly from ep.Runs).
func closureUnknownCaseAtTop(t *testing.T) {
	doc := &ast.Document{
		Cases: []ast.Case{},
		Provides: []ast.Endpoint{
			{Name: "start", Runs: []string{"missing"}},
		},
	}
	_, err := Closure(doc, "start")
	if err == nil {
		t.Fatal("expected error for unknown case referenced by endpoint Runs")
	}
}
