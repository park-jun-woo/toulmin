//ff:func feature=tangl type=analyzer control=sequence
//ff:what closureUnknownCaseNested — verifies Closure propagates an error when a nested run edge targets an unknown case
package effects

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// closureUnknownCaseNested verifies that a run edge inside a case that
// targets a nonexistent case propagates the error up through the recursive
// visit call and then through the top loop.
func closureUnknownCaseNested(t *testing.T) {
	ref := &ast.Ref{Name: "f"}
	doc := &ast.Document{
		Cases: []ast.Case{
			{
				Name: "root",
				Execs: []ast.Exec{
					{Kind: ast.DoExec, Func: ref, Node: "n1"},
					{Kind: ast.RunExec, Case: "missing", Node: "n1"},
				},
			},
		},
		Provides: []ast.Endpoint{
			{Name: "start", Runs: []string{"root"}},
		},
	}
	_, err := Closure(doc, "start")
	if err == nil {
		t.Fatal("expected error for nested run edge to unknown case")
	}
}
