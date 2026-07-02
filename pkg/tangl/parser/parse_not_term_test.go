//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseNotTerm — tests parseNotTerm for the empty-rest (condList) and non-empty-rest (compare) branches
package parser

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestParseNotTerm(t *testing.T) {
	t.Run("EmptyRestUsesCondList", func(t *testing.T) {
		it := item{
			Line: 1,
			Children: []item{
				{Text: "`x` equals 1", Line: 2},
			},
		}
		expr, err := parseNotTerm("  ", it, "test.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		cmp, ok := expr.(ast.Compare)
		if !ok {
			t.Fatalf("expected ast.Compare, got %T", expr)
		}
		if cmp.Field != "x" || cmp.Op != "equals" || cmp.Value != "1" {
			t.Errorf("unexpected compare: %+v", cmp)
		}
	})

	t.Run("EmptyRestCondListError", func(t *testing.T) {
		it := item{Line: 3}
		_, err := parseNotTerm("", it, "test.md")
		if err == nil {
			t.Fatal("expected an error from parseCondList for missing children")
		}
	})

	t.Run("NonEmptyRestUsesCompare", func(t *testing.T) {
		it := item{Line: 4}
		expr, err := parseNotTerm("`x` equals 1", it, "test.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		cmp, ok := expr.(ast.Compare)
		if !ok {
			t.Fatalf("expected ast.Compare, got %T", expr)
		}
		if cmp.Field != "x" || cmp.Op != "equals" || cmp.Value != "1" {
			t.Errorf("unexpected compare: %+v", cmp)
		}
	})

	t.Run("NonEmptyRestCompareError", func(t *testing.T) {
		it := item{Line: 5}
		_, err := parseNotTerm("notbacktick equals 1", it, "test.md")
		if err == nil {
			t.Fatal("expected an error from parseCompare for a missing backtick field")
		}
	})
}
