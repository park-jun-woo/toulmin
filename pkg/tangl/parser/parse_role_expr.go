//ff:func feature=tangl type=parser control=sequence
//ff:what parseRoleExpr — parse a ROLE_EXPR (general/counter/except rule)
package parser

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// parseRoleExpr parses "general rule", "counter rule", or "except rule".
func parseRoleExpr(s string) (ast.Role, string, bool) {
	if r, ok := takeKeyword(s, "general rule"); ok {
		return ast.GeneralRule, r, true
	}
	if r, ok := takeKeyword(s, "counter rule"); ok {
		return ast.CounterRule, r, true
	}
	if r, ok := takeKeyword(s, "except rule"); ok {
		return ast.ExceptRule, r, true
	}
	return 0, s, false
}
