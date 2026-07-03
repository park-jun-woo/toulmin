//ff:func feature=tangl type=validator control=sequence
//ff:what TestArmUndoExec — tests armUndoExec for do-arms, undo-armed, undo-unarmed, and non-do/undo branches
package validate

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestArmUndoExec(t *testing.T) {
	t.Run("DoArms", func(t *testing.T) {
		armed := map[string]bool{}
		err := armUndoExec("test.md", "c1", armed, ast.Exec{Kind: ast.DoExec, Node: "n1"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !armed["n1"] {
			t.Error("expected n1 to be armed after a do exec")
		}
	})

	t.Run("UndoArmed", func(t *testing.T) {
		armed := map[string]bool{"n1": true}
		err := armUndoExec("test.md", "c1", armed, ast.Exec{Kind: ast.UndoExec, Node: "n1"})
		if err != nil {
			t.Fatalf("unexpected error for an armed undo: %v", err)
		}
	})

	t.Run("UndoUnarmed", func(t *testing.T) {
		armed := map[string]bool{}
		err := armUndoExec("test.md", "c1", armed, ast.Exec{Kind: ast.UndoExec, Node: "n1", Line: 5})
		if err == nil {
			t.Fatal("expected an error for an unarmed undo")
		}
		if !strings.Contains(err.Error(), "has no preceding do") {
			t.Errorf("expected preceding-do error, got %v", err)
		}
	})

	t.Run("OtherKindIsNoOp", func(t *testing.T) {
		armed := map[string]bool{}
		err := armUndoExec("test.md", "c1", armed, ast.Exec{Kind: ast.RunExec, Node: "n1"})
		if err != nil {
			t.Fatalf("unexpected error for a run exec: %v", err)
		}
		if len(armed) != 0 {
			t.Errorf("expected armed to stay empty, got %+v", armed)
		}
	})
}
