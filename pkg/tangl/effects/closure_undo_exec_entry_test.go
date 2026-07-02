//ff:func feature=tangl type=analyzer control=sequence
//ff:what closureUndoExecEntry — exercises the UndoExec switch branch directly, independent of the transfer fixture
package effects

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// closureUndoExecEntry exercises the UndoExec switch branch directly
// (independent of the transfer fixture) to keep it self-contained.
func closureUndoExecEntry(t *testing.T) {
	ref := &ast.Ref{Name: "f"}
	doc := &ast.Document{
		Cases: []ast.Case{
			{
				Name:  "root",
				Execs: []ast.Exec{{Kind: ast.UndoExec, Func: ref, Node: "n1"}},
			},
		},
		Provides: []ast.Endpoint{
			{Name: "start", Runs: []string{"root"}},
		},
	}
	entries, err := Closure(doc, "start")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 || entries[0].Kind != "undo" {
		t.Fatalf("expected one undo entry, got %+v", entries)
	}
}
