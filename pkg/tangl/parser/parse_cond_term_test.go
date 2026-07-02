//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseCondTerm — tests parseCondTerm for either success/error, not success/error, and compare success/error branches
package parser

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestParseCondTerm(t *testing.T) {
	t.Run("EitherSuccess", func(t *testing.T) {
		it := item{
			Text: "either",
			Children: []item{
				{Text: "`amount` is empty", Line: 1},
			},
		}
		expr, err := parseCondTerm(it, "test.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if _, ok := expr.(ast.Either); !ok {
			t.Fatalf("expected ast.Either, got %T", expr)
		}
	})

	t.Run("EitherError", func(t *testing.T) {
		it := item{Text: "either"}
		_, err := parseCondTerm(it, "test.md")
		if err == nil {
			t.Fatal("expected an error for an 'either' with no children")
		}
	})

	t.Run("NotSuccess", func(t *testing.T) {
		it := item{Text: "not `amount` is empty", Line: 1}
		expr, err := parseCondTerm(it, "test.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if _, ok := expr.(ast.Not); !ok {
			t.Fatalf("expected ast.Not, got %T", expr)
		}
	})

	t.Run("NotError", func(t *testing.T) {
		it := item{Text: "not bogus", Line: 1}
		_, err := parseCondTerm(it, "test.md")
		if err == nil {
			t.Fatal("expected an error for a malformed 'not' term")
		}
	})

	t.Run("CompareSuccess", func(t *testing.T) {
		it := item{Text: "`amount` is empty", Line: 1}
		expr, err := parseCondTerm(it, "test.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		cmp, ok := expr.(ast.Compare)
		if !ok || cmp.Field != "amount" {
			t.Fatalf("expected Compare{amount}, got %+v", expr)
		}
	})

	t.Run("CompareError", func(t *testing.T) {
		it := item{Text: "bogus", Line: 1}
		_, err := parseCondTerm(it, "test.md")
		if err == nil {
			t.Fatal("expected an error for a malformed compare term")
		}
	})
}
