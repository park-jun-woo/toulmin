//ff:func feature=tangl type=codegen control=iteration dimension=1 pattern=recursive
//ff:what collectRuleFlags — scans every Rules condition tree to set the rule-helper flags
package gen

import "github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"

// collectRuleFlags walks every tangl:Rules condition tree with an active
// exists-binder scope and sets the Needs* flags for each construct the
// closure-composition codegen will emit: comparison/ctx/spec/walk value
// helpers, the and/or/not/any combinators, the holds passthrough, and
// strconv.Itoa for binder-indexed Where strings. NeedsCtxVal and NeedsSpecVal
// imply NeedsWalk because their emitted helpers call tanglWalk.
func collectRuleFlags(rules []ast.InlineRule, f *genFlags) {
	operand := func(o ast.Operand, active map[string]bool) {
		switch v := o.(type) {
		case ast.Path:
			if len(v.Steps) > 0 && active[v.Steps[0]] {
				f.NeedsItoa = true
				if len(v.Steps) > 1 {
					f.NeedsWalk = true
				}
			} else {
				f.NeedsCtxVal = true
			}
		case ast.SpecRef:
			f.NeedsSpecVal = true
		}
	}
	var visit func(e ast.Expr, active map[string]bool)
	visit = func(e ast.Expr, active map[string]bool) {
		switch v := e.(type) {
		case ast.Compare:
			f.NeedsCompareVal = true
			if v.Lhs == nil {
				f.NeedsCtxVal = true
			} else {
				operand(v.Lhs, active)
			}
			if v.Rhs != nil {
				operand(v.Rhs, active)
			}
		case ast.Logic:
			if v.Op == "or" {
				f.NeedsOr = true
			} else {
				f.NeedsAnd = true
			}
			for _, t := range v.Terms {
				visit(t, active)
			}
		case ast.Either:
			f.NeedsOr = true
			for _, t := range v.Terms {
				visit(t, active)
			}
		case ast.Not:
			f.NeedsNot = true
			visit(v.Term, active)
		case ast.Exists:
			f.NeedsAny = true
			operand(v.Source, active)
			inner := make(map[string]bool, len(active)+1)
			for k := range active {
				inner[k] = true
			}
			inner[v.Binder] = true
			visit(v.Cond, inner)
		case ast.Call:
			f.NeedsHolds = true
		}
	}
	for _, r := range rules {
		visit(r.Cond, map[string]bool{})
	}
	if f.NeedsCtxVal || f.NeedsSpecVal {
		f.NeedsWalk = true
	}
}
