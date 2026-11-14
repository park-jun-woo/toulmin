//ff:func feature=tangl type=parser control=sequence
//ff:what TestApplyNodeClauseQualified — applyNodeClause qualified clause without using/checking, after checking, range error, duplicate, and bad level
package parser

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestApplyNodeClauseQualified(t *testing.T) {
	t.Run("QualifiedOnlyWithLevel", func(t *testing.T) {
		n := &ast.Node{}
		if err := applyNodeClause(n, "qualified 0.8, level fail", 8, "test.md"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if n.Qualified == nil || *n.Qualified != 0.8 || n.Level != ast.LevelFail || n.Using != nil || n.Checking != nil {
			t.Fatalf("got %+v, want qualified 0.8 level fail without using/checking", n)
		}
	})

	t.Run("CheckingQualified", func(t *testing.T) {
		n := &ast.Node{}
		if err := applyNodeClause(n, "checking `inner` qualified 0.5", 9, "test.md"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if n.Checking == nil || n.Checking.Name != "inner" || n.Qualified == nil || *n.Qualified != 0.5 {
			t.Fatalf("got %+v, want checking inner qualified 0.5", n)
		}
	})

	t.Run("QualifiedOutOfRange", func(t *testing.T) {
		n := &ast.Node{}
		err := applyNodeClause(n, "qualified 1.5", 10, "test.md")
		if err == nil || !strings.Contains(err.Error(), "out of range") {
			t.Fatalf("expected an out-of-range error, got %v", err)
		}
	})

	t.Run("SecondQualifiedIsTrailing", func(t *testing.T) {
		n := &ast.Node{}
		err := applyNodeClause(n, "using `fn` qualified 0.5 qualified 0.7", 11, "test.md")
		if err == nil || !strings.Contains(err.Error(), "unexpected trailing text") {
			t.Fatalf("expected a trailing-text error, got %v", err)
		}
	})

	t.Run("BadLevel", func(t *testing.T) {
		n := &ast.Node{}
		if err := applyNodeClause(n, "qualified 0.5, level loud", 12, "test.md"); err == nil {
			t.Fatal("expected an error for an unknown level")
		}
	})
}
