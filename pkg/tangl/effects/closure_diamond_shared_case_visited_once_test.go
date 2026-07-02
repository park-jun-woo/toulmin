//ff:func feature=tangl type=analyzer control=sequence
//ff:what closureDiamondSharedCaseVisitedOnce — verifies a case reachable via two run paths executes only once
package effects

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// closureDiamondSharedCaseVisitedOnce verifies that a case reachable via two
// independent run paths (a diamond) is only executed once: the second visit
// hits state[name]==2 (done) and is skipped.
func closureDiamondSharedCaseVisitedOnce(t *testing.T) {
	ref := &ast.Ref{Name: "shared"}
	doc := &ast.Document{
		Cases: []ast.Case{
			{
				Name:  "top1",
				Execs: []ast.Exec{{Kind: ast.RunExec, Case: "shared", Node: "n1"}},
			},
			{
				Name:  "top2",
				Execs: []ast.Exec{{Kind: ast.RunExec, Case: "shared", Node: "n1"}},
			},
			{
				Name:  "shared",
				Execs: []ast.Exec{{Kind: ast.DoExec, Func: ref, Node: "n1"}},
			},
		},
		Provides: []ast.Endpoint{
			{Name: "start", Runs: []string{"top1", "top2"}},
		},
	}
	entries, err := Closure(doc, "start")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected shared case's do entry exactly once, got %+v", entries)
	}
}
