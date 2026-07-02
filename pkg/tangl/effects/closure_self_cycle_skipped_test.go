//ff:func feature=tangl type=analyzer control=sequence
//ff:what closureSelfCycleSkipped — verifies Closure does not infinite-loop when a case runs itself
package effects

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// closureSelfCycleSkipped verifies that a case that runs itself does not
// infinite-loop: the state[name]==1 (visiting) branch short-circuits the
// recursive re-entry and Closure returns without error.
func closureSelfCycleSkipped(t *testing.T) {
	ref := &ast.Ref{Name: "f"}
	doc := &ast.Document{
		Cases: []ast.Case{
			{
				Name: "cyclic",
				Execs: []ast.Exec{
					{Kind: ast.RunExec, Case: "cyclic", Node: "n1"},
					{Kind: ast.DoExec, Func: ref, Node: "n1"},
				},
			},
		},
		Provides: []ast.Endpoint{
			{Name: "start", Runs: []string{"cyclic"}},
		},
	}
	entries, err := Closure(doc, "start")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected exactly one entry (self-run skipped), got %+v", entries)
	}
}
