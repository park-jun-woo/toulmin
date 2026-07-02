//ff:func feature=tangl type=validator control=sequence
//ff:what TestCheckUndoRequiresDo — tests checkUndoRequiresDo for no-cases, armed-undo, unarmed-undo, and non-do/undo-skip branches
package validate

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestCheckUndoRequiresDo(t *testing.T) {
	t.Run("NoCases", func(t *testing.T) {
		doc := &ast.Document{}
		errs := checkUndoRequiresDo(doc)
		if len(errs) != 0 {
			t.Fatalf("expected no errors, got %v", errs)
		}
	})

	t.Run("ArmedUndo", func(t *testing.T) {
		doc := &ast.Document{
			Cases: []ast.Case{
				{
					Name: "x",
					Execs: []ast.Exec{
						{Kind: ast.DoExec, Node: "n1"},
						{Kind: ast.UndoExec, Node: "n1"},
					},
				},
			},
		}
		errs := checkUndoRequiresDo(doc)
		if len(errs) != 0 {
			t.Fatalf("expected no errors, got %v", errs)
		}
	})

	t.Run("UnarmedUndo", func(t *testing.T) {
		doc := &ast.Document{
			Cases: []ast.Case{
				{
					Name: "x",
					Execs: []ast.Exec{
						{Kind: ast.UndoExec, Node: "n1", Line: 5},
					},
				},
			},
		}
		errs := checkUndoRequiresDo(doc)
		if len(errs) != 1 {
			t.Fatalf("expected 1 error, got %v", errs)
		}
		if !strings.Contains(errs[0].Error(), "has no preceding do") {
			t.Errorf("expected preceding do error, got %v", errs[0])
		}
	})

	t.Run("SkipsRunExec", func(t *testing.T) {
		doc := &ast.Document{
			Cases: []ast.Case{
				{
					Name: "x",
					Execs: []ast.Exec{
						{Kind: ast.RunExec, Node: "n1", Case: "other"},
					},
				},
			},
		}
		errs := checkUndoRequiresDo(doc)
		if len(errs) != 0 {
			t.Fatalf("expected no errors for run exec, got %v", errs)
		}
	})
}
