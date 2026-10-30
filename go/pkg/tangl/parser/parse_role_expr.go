//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what parseRoleExpr — parse a ROLE_EXPR (general/counter/except rule)
package parser

import (
	"strings"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// parseRoleExpr parses "general rule", "counter rule", or "except rule".
func parseRoleExpr(s string) (ast.Role, string, bool) {
	for _, candidate := range []struct {
		text string
		role ast.Role
	}{{"general rule", ast.GeneralRule}, {"counter rule", ast.CounterRule}, {"except rule", ast.ExceptRule}} {
		if strings.HasPrefix(strings.TrimSpace(s), candidate.text+",") {
			trimmed := strings.TrimSpace(s)
			return candidate.role, trimmed[len(candidate.text):], true
		}
	}
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
