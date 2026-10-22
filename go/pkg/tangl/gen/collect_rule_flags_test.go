//ff:func feature=tangl type=codegen control=sequence
//ff:what TestCollectRuleFlags — tests the Needs* flag scan across every rule construct
package gen

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestCollectRuleFlags(t *testing.T) {
	t.Run("flat compare", func(t *testing.T) {
		var f genFlags
		rules := []ast.InlineRule{{Name: "r", Cond: ast.Compare{Field: "x", Op: "equals", Value: "1"}}}
		collectRuleFlags(rules, &f)
		if !f.NeedsCompareVal || !f.NeedsCtxVal || !f.NeedsWalk {
			t.Errorf("flat compare should need CompareVal/CtxVal/Walk: %+v", f)
		}
	})
	t.Run("call only skips walk", func(t *testing.T) {
		var f genFlags
		collectRuleFlags([]ast.InlineRule{{Name: "r", Cond: ast.Call{Fn: ast.Ref{Name: "p"}}}}, &f)
		if !f.NeedsHolds || f.NeedsWalk || f.NeedsCtxVal {
			t.Errorf("a lone call should only need Holds: %+v", f)
		}
	})
	t.Run("combinators", func(t *testing.T) {
		var f genFlags
		cmp := ast.Compare{Field: "x", Op: "is empty"}
		rules := []ast.InlineRule{{Name: "r", Cond: ast.Logic{Op: "and", Terms: []ast.Expr{
			ast.Logic{Op: "or", Terms: []ast.Expr{cmp}},
			ast.Either{Terms: []ast.Expr{cmp}},
			ast.Not{Term: cmp},
		}}}}
		collectRuleFlags(rules, &f)
		if !f.NeedsAnd || !f.NeedsOr || !f.NeedsNot || !f.NeedsCompareVal {
			t.Errorf("combinators should need And/Or/Not/CompareVal: %+v", f)
		}
	})
	t.Run("exists operand forms", func(t *testing.T) {
		var f genFlags
		rules := []ast.InlineRule{{Name: "r", Cond: ast.Exists{
			Source: ast.Path{Steps: []string{"coupons"}},
			Binder: "c",
			Cond: ast.Logic{Op: "and", Terms: []ast.Expr{
				ast.Compare{Lhs: ast.Path{Steps: []string{"c"}}, Op: "is not empty"},
				ast.Compare{Lhs: ast.Path{Steps: []string{"c", "minPrice"}}, Op: "is at most", Rhs: ast.SpecRef{Steps: []string{"d", "max"}}},
				ast.Compare{Lhs: ast.Path{Steps: []string{"other"}}, Op: "equals", Rhs: ast.Path{}},
				ast.Exists{
					Source: ast.Path{Steps: []string{"c", "items"}},
					Binder: "d",
					Cond:   ast.Compare{Lhs: ast.Path{Steps: []string{"d"}}, Op: "is not empty"},
				},
			}},
		}}}
		collectRuleFlags(rules, &f)
		if !f.NeedsAny || !f.NeedsItoa || !f.NeedsWalk || !f.NeedsSpecVal || !f.NeedsCtxVal {
			t.Errorf("exists forms should need Any/Itoa/Walk/SpecVal/CtxVal: %+v", f)
		}
	})
}
