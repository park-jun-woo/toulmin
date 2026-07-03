//ff:func feature=tangl type=codegen control=selection
//ff:what levelArg — maps a Node's Level to its toulmin.Level codegen argument
package gen

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// levelArg maps a Node's Level to the toulmin.Level identifier it codegens
// to: fail maps to toulmin.Fail, review maps to toulmin.Review. LevelUnset
// returns "" — the caller filters it out before emitting a .Level(...) call.
func levelArg(l ast.Level) string {
	switch l {
	case ast.LevelFail:
		return "toulmin.Fail"
	case ast.LevelReview:
		return "toulmin.Review"
	default:
		return ""
	}
}
