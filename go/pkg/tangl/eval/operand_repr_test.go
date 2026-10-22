//ff:func feature=tangl type=engine control=sequence
//ff:what TestReprOperand — tests Where rendering for Path, SpecRef, Call, and binder-indexed paths
package eval

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/go/pkg/toulmin"
)

func TestReprOperand(t *testing.T) {
	sc := newScope(toulmin.NewContext(), nil)

	t.Run("Path", func(t *testing.T) {
		if got := reprOperand(ast.Path{Steps: []string{"author", "trustScore"}}, sc); got != "author.trustScore" {
			t.Fatalf("got %q", got)
		}
	})
	t.Run("SpecRef", func(t *testing.T) {
		if got := reprOperand(ast.SpecRef{Steps: []string{"discount", "minPrice"}}, sc); got != "discount.minPrice" {
			t.Fatalf("got %q", got)
		}
	})
	t.Run("CallBare", func(t *testing.T) {
		if got := reprOperand(ast.Call{Fn: ast.Ref{Name: "score"}}, sc); got != "score" {
			t.Fatalf("got %q", got)
		}
	})
	t.Run("CallAliased", func(t *testing.T) {
		if got := reprOperand(ast.Call{Fn: ast.Ref{Alias: "nlp", Name: "score"}}, sc); got != "nlp.score" {
			t.Fatalf("got %q", got)
		}
	})
	t.Run("UnknownOperand", func(t *testing.T) {
		if got := reprOperand(nil, sc); got != "" {
			t.Fatalf("nil operand should render empty, got %q", got)
		}
	})
	t.Run("BinderIndexed", func(t *testing.T) {
		sc.binders["c"] = &binding{source: ast.Path{Steps: []string{"coupons"}}, index: 2}
		if got := reprOperand(ast.Path{Steps: []string{"c", "minPrice"}}, sc); got != "coupons[2].minPrice" {
			t.Fatalf("got %q", got)
		}
		if got := reprOperand(ast.Path{Steps: []string{"c"}}, sc); got != "coupons[2]" {
			t.Fatalf("bare binder got %q", got)
		}
	})
}
