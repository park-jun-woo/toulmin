//ff:func feature=tangl type=codegen control=iteration dimension=1 pattern=recursive
//ff:what collectRuleCallAliases — collects See aliases referenced by host Calls inside Rules
package gen

import "github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"

// collectRuleCallAliases walks every tangl:Rules condition tree and records
// into used the package alias of every alias-qualified host Call — both
// predicate-position Calls (`alias.f holds`) and value-position Calls
// (`result of alias.f` in Compare operands or an Exists source) — so
// buildImports emits the See import the generated call sites need.
func collectRuleCallAliases(rules []ast.InlineRule, used map[string]bool) {
	operand := func(o ast.Operand) {
		if c, ok := o.(ast.Call); ok && c.Fn.Alias != "" {
			used[c.Fn.Alias] = true
		}
	}
	var visit func(e ast.Expr)
	visit = func(e ast.Expr) {
		switch v := e.(type) {
		case ast.Compare:
			if v.Lhs != nil {
				operand(v.Lhs)
			}
			if v.Rhs != nil {
				operand(v.Rhs)
			}
		case ast.Logic:
			for _, t := range v.Terms {
				visit(t)
			}
		case ast.Either:
			for _, t := range v.Terms {
				visit(t)
			}
		case ast.Not:
			visit(v.Term)
		case ast.Exists:
			operand(v.Source)
			visit(v.Cond)
		case ast.Call:
			if v.Fn.Alias != "" {
				used[v.Fn.Alias] = true
			}
		}
	}
	for _, r := range rules {
		visit(r.Cond)
	}
}
