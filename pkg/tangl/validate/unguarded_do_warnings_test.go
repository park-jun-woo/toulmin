//ff:func feature=tangl type=validator control=sequence
//ff:what TestUnguardedDoWarnings — tests unguardedDoWarnings for unreached-case-skip, non-do-skip, once-guarded-skip, and unguarded-warning branches
package validate

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestUnguardedDoWarnings(t *testing.T) {
	t.Run("UnreachedCase", func(t *testing.T) {
		doc := &ast.Document{
			Cases: []ast.Case{
				{Name: "x", Execs: []ast.Exec{
					{Kind: ast.DoExec, Once: false},
				}},
			},
		}
		reached := map[string]bool{}
		warnings := unguardedDoWarnings(doc, reached)
		if len(warnings) != 0 {
			t.Fatalf("expected no warnings for unreached case, got %v", warnings)
		}
	})

	t.Run("SkipsNonDoExec", func(t *testing.T) {
		doc := &ast.Document{
			Cases: []ast.Case{
				{Name: "x", Execs: []ast.Exec{
					{Kind: ast.UndoExec, Once: false},
				}},
			},
		}
		reached := map[string]bool{"x": true}
		warnings := unguardedDoWarnings(doc, reached)
		if len(warnings) != 0 {
			t.Fatalf("expected no warnings for non-do exec, got %v", warnings)
		}
	})

	t.Run("SkipsOnceGuarded", func(t *testing.T) {
		doc := &ast.Document{
			Cases: []ast.Case{
				{Name: "x", Execs: []ast.Exec{
					{Kind: ast.DoExec, Once: true},
				}},
			},
		}
		reached := map[string]bool{"x": true}
		warnings := unguardedDoWarnings(doc, reached)
		if len(warnings) != 0 {
			t.Fatalf("expected no warnings for once-guarded do, got %v", warnings)
		}
	})

	t.Run("WarnsOnUnguardedDo", func(t *testing.T) {
		doc := &ast.Document{
			Path: "doc.md",
			Cases: []ast.Case{
				{Name: "x", Execs: []ast.Exec{
					{Kind: ast.DoExec, Once: false, Line: 7, Func: &ast.Ref{Name: "fn"}},
				}},
			},
		}
		reached := map[string]bool{"x": true}
		warnings := unguardedDoWarnings(doc, reached)
		if len(warnings) != 1 {
			t.Fatalf("expected 1 warning, got %v", warnings)
		}
		if !strings.Contains(warnings[0], "no once guard") {
			t.Errorf("expected no once guard warning, got %v", warnings[0])
		}
	})
}
