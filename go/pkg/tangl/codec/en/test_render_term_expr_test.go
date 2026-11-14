//ff:func feature=tangl type=parser control=sequence pattern=recursive
//ff:what TestRenderTermExprTerms — renderTerm covers the exists and host-call arms (one-line and nested)
package en

import (
	"reflect"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// TestRenderTermExprTerms exercises the exists and call arms: a host predicate
// Call (aliased and bare, one line, with and without a chaining prefix), a lone
// exists (one line via renderExists), and a nested exists (head plus list).
func TestRenderTermExprTerms(t *testing.T) {
	// Aliased Call term, no prefix.
	got := renderTerm("", ast.Call{Fn: ast.Ref{Alias: "check", Name: "isSpam"}}, 0)
	want := []string{"- `check`.`isSpam` holds"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Call aliased: got %v, want %v", got, want)
	}

	// Bare Call term, "or" prefix.
	got = renderTerm("or", ast.Call{Fn: ast.Ref{Name: "manualFlag"}}, 0)
	want = []string{"- or `manualFlag` holds"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Call bare/prefixed: got %v, want %v", got, want)
	}

	// Lone-compare exists, "and" prefix: one line.
	lone := ast.Exists{Source: ast.Path{Steps: []string{"coupons"}}, Binder: "c",
		Cond: ast.Compare{Field: "c", Op: "is at most", Rhs: ast.Path{Steps: []string{"basePrice"}}}}
	got = renderTerm("and", lone, 0)
	want = []string{"- and any of `coupons` as `c` satisfies `c` is at most `basePrice`"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Exists lone: got %v, want %v", got, want)
	}

	// Nested exists: head line plus a depth+1 list.
	nested := ast.Exists{Source: ast.SpecRef{Steps: []string{"promotions"}}, Binder: "p",
		Cond: ast.Logic{Op: "and", Terms: []ast.Expr{
			ast.Compare{Op: "equals", Value: "true", Lhs: ast.Path{Steps: []string{"p", "active"}}},
			ast.Compare{Op: "is at least", Value: "10", Lhs: ast.Path{Steps: []string{"p", "discount"}}},
		}}}
	got = renderTerm("", nested, 0)
	want = []string{
		"- any of spec `promotions` as `p` satisfies",
		"  - `p`.`active` equals true",
		"  - and `p`.`discount` is at least 10",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Exists nested: got %v, want %v", got, want)
	}
}
