//ff:func feature=tangl type=parser control=sequence
//ff:what TestApplyNodeClause_Level — tests applyNodeClause level parsing for using/checking clauses, bad tokens, and the no-clause fallback
package parser

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestApplyNodeClause_Level(t *testing.T) {
	t.Run("UsingLevelFail", func(t *testing.T) {
		n := &ast.Node{}
		err := applyNodeClause(n, "using `fn`, level fail", 8, "test.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if n.Using == nil || n.Using.Name != "fn" {
			t.Fatalf("expected Using to be set to fn, got %+v", n.Using)
		}
		if n.Level != ast.LevelFail {
			t.Fatalf("expected Level = LevelFail, got %v", n.Level)
		}
	})

	t.Run("UsingLevelReview", func(t *testing.T) {
		n := &ast.Node{}
		err := applyNodeClause(n, "using `fn`, level review", 9, "test.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if n.Level != ast.LevelReview {
			t.Fatalf("expected Level = LevelReview, got %v", n.Level)
		}
	})

	t.Run("CheckingLevelFail", func(t *testing.T) {
		n := &ast.Node{}
		err := applyNodeClause(n, "checking `otherCase`, level fail", 10, "test.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if n.Checking == nil || n.Checking.Name != "otherCase" || n.Level != ast.LevelFail {
			t.Fatalf("expected Checking=otherCase Level=LevelFail, got %+v / %v", n.Checking, n.Level)
		}
	})

	t.Run("BadLevelToken", func(t *testing.T) {
		n := &ast.Node{}
		err := applyNodeClause(n, "using `fn`, level foo", 11, "test.md")
		if err == nil || !strings.Contains(err.Error(), "expected 'fail' or 'review'") {
			t.Fatalf("expected bad-level-token error, got %v", err)
		}
	})

	t.Run("NoLevelUnaffected", func(t *testing.T) {
		n := &ast.Node{}
		err := applyNodeClause(n, "using `fn`", 12, "test.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if n.Level != ast.LevelUnset {
			t.Fatalf("expected Level = LevelUnset when no clause, got %v", n.Level)
		}
	})
}
