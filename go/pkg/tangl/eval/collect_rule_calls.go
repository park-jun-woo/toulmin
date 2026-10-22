//ff:func feature=tangl type=engine control=selection pattern=recursive
//ff:what collectRuleCalls — collects bare predicate Call names that are document rules
package eval

import "github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"

// collectRuleCalls appends to out the names of every bare (alias-less)
// predicate-position Call in e that names a document rule (present in ruleSet).
// Value-position operands are not traversed: doc-local rules are predicates,
// not value functions, so only predicate edges count.
func collectRuleCalls(e ast.Expr, ruleSet map[string]bool, out *[]string) {
	switch v := e.(type) {
	case ast.Logic:
		for _, t := range v.Terms {
			collectRuleCalls(t, ruleSet, out)
		}
	case ast.Either:
		for _, t := range v.Terms {
			collectRuleCalls(t, ruleSet, out)
		}
	case ast.Not:
		collectRuleCalls(v.Term, ruleSet, out)
	case ast.Exists:
		collectRuleCalls(v.Cond, ruleSet, out)
	case ast.Call:
		if v.Fn.Alias == "" && ruleSet[v.Fn.Name] {
			*out = append(*out, v.Fn.Name)
		}
	}
}
