//ff:func feature=tangl type=parser control=sequence
//ff:what TestRenderExists — renderExists covers lone-compare (one line), nested body, and the term vs rule-head lead
package en

import (
	"reflect"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// TestRenderExists exercises a lone-compare body (single line following
// "satisfies"), a nested and/or body (head line plus a depth+1 list), and both
// leads — the "- " term marker and the "`name` when " rule head — over the same
// helper so the two dispatch sites stay in lock-step.
func TestRenderExists(t *testing.T) {
	cmp := ast.Compare{Field: "c", Op: "is at most", Rhs: ast.Path{Steps: []string{"basePrice"}}}
	lone := ast.Exists{Source: ast.Path{Steps: []string{"coupons"}}, Binder: "c", Cond: cmp}

	// Term position, lone-compare body: one line.
	got := renderExists("- ", lone, 0)
	want := []string{"- any of `coupons` as `c` satisfies `c` is at most `basePrice`"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("term lone: got %v, want %v", got, want)
	}

	// Rule-head position, lone-compare body: one line.
	got = renderExists("`hasCoupon` when ", lone, 0)
	want = []string{"`hasCoupon` when any of `coupons` as `c` satisfies `c` is at most `basePrice`"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("head lone: got %v, want %v", got, want)
	}

	// Nested and-chain body: head line plus a depth+1 nested list.
	nested := ast.Exists{Source: ast.Path{Steps: []string{"promotions"}}, Binder: "p", Cond: ast.Logic{Op: "and", Terms: []ast.Expr{
		ast.Compare{Field: "p", Op: "equals", Value: "true", Lhs: ast.Path{Steps: []string{"p", "active"}}},
		ast.Compare{Field: "p", Op: "is at least", Value: "10", Lhs: ast.Path{Steps: []string{"p", "discount"}}},
	}}}
	got = renderExists("- ", nested, 0)
	want = []string{
		"- any of `promotions` as `p` satisfies",
		"  - `p`.`active` equals true",
		"  - and `p`.`discount` is at least 10",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("nested: got %v, want %v", got, want)
	}
}
