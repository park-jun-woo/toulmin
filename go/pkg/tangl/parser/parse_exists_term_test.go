//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseExistsTerm — tests parseExistsTerm for one-line and nested success branches
package parser

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestParseExistsTerm(t *testing.T) {
	t.Run("OneLine", func(t *testing.T) {
		it := item{Text: "any of `coupons` as `c` satisfies `c`.`minPrice` is at most `basePrice`", Line: 1}
		expr, err := parseExistsTerm(it.Text, it, "test.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		ex, ok := expr.(ast.Exists)
		if !ok {
			t.Fatalf("expected ast.Exists, got %T", expr)
		}
		if ex.Binder != "c" {
			t.Errorf("expected binder c, got %q", ex.Binder)
		}
		src, ok := ex.Source.(ast.Path)
		if !ok || len(src.Steps) != 1 || src.Steps[0] != "coupons" {
			t.Errorf("expected source Path[coupons], got %+v", ex.Source)
		}
		cmp, ok := ex.Cond.(ast.Compare)
		if !ok {
			t.Fatalf("expected Compare cond, got %T", ex.Cond)
		}
		p, ok := cmp.Lhs.(ast.Path)
		if !ok || p.Steps[0] != "c" || p.Steps[1] != "minPrice" {
			t.Errorf("expected cond Lhs c.minPrice, got %+v", cmp.Lhs)
		}
	})

	t.Run("Nested", func(t *testing.T) {
		it := item{
			Text: "any of `coupons` as `c` satisfies",
			Line: 2,
			Children: []item{
				{Text: "`c`.`active` equals true", Line: 3},
				{Text: "and `c`.`minPrice` is at most `basePrice`", Line: 4},
			},
		}
		expr, err := parseExistsTerm(it.Text, it, "test.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		ex, ok := expr.(ast.Exists)
		if !ok {
			t.Fatalf("expected ast.Exists, got %T", expr)
		}
		if _, ok := ex.Cond.(ast.Logic); !ok {
			t.Fatalf("expected Logic cond from and-chained children, got %T", ex.Cond)
		}
	})
}
