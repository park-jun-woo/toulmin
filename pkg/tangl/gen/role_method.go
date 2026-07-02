//ff:func feature=tangl type=codegen control=selection
//ff:what roleMethod — maps a Node's Role to its toulmin.Graph registration method
package gen

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// roleMethod maps a Node's Role to the toulmin.Graph registration method
// it codegens to: general maps to Rule, counter maps to Counter, except
// maps to Except.
func roleMethod(r ast.Role) string {
	switch r {
	case ast.CounterRule:
		return "Counter"
	case ast.ExceptRule:
		return "Except"
	default:
		return "Rule"
	}
}
