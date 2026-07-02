//ff:func feature=tangl type=parser control=sequence
//ff:what TestApplyNodeClause — tests applyNodeClause for using/checking success, error, and trailing-text branches, plus the neither-clause fallback
package parser

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestApplyNodeClause(t *testing.T) {
	t.Run("UsingSuccess", func(t *testing.T) {
		n := &ast.Node{}
		err := applyNodeClause(n, "using `fn`", 1, "test.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if n.Using == nil || n.Using.Name != "fn" {
			t.Fatalf("expected Using to be set to fn, got %+v", n.Using)
		}
	})

	t.Run("UsingError", func(t *testing.T) {
		n := &ast.Node{}
		err := applyNodeClause(n, "using notabacktick", 2, "test.md")
		if err == nil {
			t.Fatal("expected an error for a missing function reference")
		}
		if !strings.Contains(err.Error(), "expected function reference after 'using'") {
			t.Errorf("expected function reference error, got %v", err)
		}
	})

	t.Run("UsingTrailingText", func(t *testing.T) {
		n := &ast.Node{}
		err := applyNodeClause(n, "using `fn` extra", 3, "test.md")
		if err == nil {
			t.Fatal("expected an error for trailing text after using clause")
		}
		if !strings.Contains(err.Error(), "unexpected trailing text") {
			t.Errorf("expected trailing text error, got %v", err)
		}
	})

	t.Run("CheckingSuccess", func(t *testing.T) {
		n := &ast.Node{}
		err := applyNodeClause(n, "checking `otherCase`", 4, "test.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if n.Checking != "otherCase" {
			t.Fatalf("expected Checking to be otherCase, got %q", n.Checking)
		}
	})

	t.Run("CheckingError", func(t *testing.T) {
		n := &ast.Node{}
		err := applyNodeClause(n, "checking notabacktick", 5, "test.md")
		if err == nil {
			t.Fatal("expected an error for a missing case name")
		}
		if !strings.Contains(err.Error(), "expected backtick-quoted case name after 'checking'") {
			t.Errorf("expected case name error, got %v", err)
		}
	})

	t.Run("CheckingTrailingText", func(t *testing.T) {
		n := &ast.Node{}
		err := applyNodeClause(n, "checking `otherCase` extra", 6, "test.md")
		if err == nil {
			t.Fatal("expected an error for trailing text after checking clause")
		}
		if !strings.Contains(err.Error(), "unexpected trailing text") {
			t.Errorf("expected trailing text error, got %v", err)
		}
	})

	t.Run("NeitherClause", func(t *testing.T) {
		n := &ast.Node{}
		err := applyNodeClause(n, "something else", 7, "test.md")
		if err == nil {
			t.Fatal("expected an error when neither using nor checking is present")
		}
		if !strings.Contains(err.Error(), "unexpected trailing text") {
			t.Errorf("expected trailing text error, got %v", err)
		}
	})
}
