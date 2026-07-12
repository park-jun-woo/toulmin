//ff:func feature=tangl type=parser control=sequence pattern=recursive
//ff:what TestRenderRuleExprHeads — renderRule covers the one-line lone-Exists/lone-Call rule heads
package en

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// TestRenderRuleExprHeads checks the one-line lone-Exists and lone-Call rule
// heads (§design-6 decision (a)): a lone-compare-body exists is fully one line, a
// nested-body exists is a head plus list, and a predicate Call is one line.
func TestRenderRuleExprHeads(t *testing.T) {
	// Lone Call rule head.
	call := renderRule(ast.InlineRule{Name: "spammy", Cond: ast.Call{Fn: ast.Ref{Alias: "check", Name: "isSpam"}}})
	if call != "`spammy` when `check`.`isSpam` holds" {
		t.Errorf("call head = %q", call)
	}

	// Lone Exists with a lone-compare body: fully one line.
	lone := renderRule(ast.InlineRule{Name: "hasCoupon", Cond: ast.Exists{
		Source: ast.Path{Steps: []string{"coupons"}}, Binder: "c",
		Cond: ast.Compare{Field: "c", Op: "is at most", Rhs: ast.Path{Steps: []string{"basePrice"}}},
	}})
	if lone != "`hasCoupon` when any of `coupons` as `c` satisfies `c` is at most `basePrice`" {
		t.Errorf("exists lone head = %q", lone)
	}

	// Exists with a nested body: head line plus a depth-1 list.
	nested := renderRule(ast.InlineRule{Name: "activePromotion", Cond: ast.Exists{
		Source: ast.Path{Steps: []string{"promotions"}}, Binder: "p",
		Cond: ast.Logic{Op: "and", Terms: []ast.Expr{
			ast.Compare{Op: "equals", Value: "true", Lhs: ast.Path{Steps: []string{"p", "active"}}},
			ast.Compare{Op: "is at least", Value: "10", Lhs: ast.Path{Steps: []string{"p", "discount"}}},
		}},
	}})
	wantNested := "`activePromotion` when any of `promotions` as `p` satisfies\n" +
		"  - `p`.`active` equals true\n" +
		"  - and `p`.`discount` is at least 10"
	if nested != wantNested {
		t.Errorf("exists nested head =\n%q\nwant\n%q", nested, wantNested)
	}
}
