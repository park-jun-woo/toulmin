//ff:func feature=tangl type=codegen control=sequence
//ff:what renderClosureCall — emits a predicate-position Call as a tanglHolds closure literal
package gen

import (
	"strconv"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// renderClosureCall emits a predicate-position Call (`f holds`) as a
// func() (bool, any) closure literal: it invokes the resolved Go symbol
// (alias.Name verbatim, or the goIdent local identifier for a bare name) and
// routes the result through tanglHolds, which passes a non-nil host evidence
// through unchanged — matching pkg/tangl/eval.compileCall.
func renderClosureCall(c ast.Call) string {
	repr := c.Fn.Name
	if c.Fn.Alias != "" {
		repr = c.Fn.Alias + "." + c.Fn.Name
	}
	return "func() (bool, any) {\nb, ev := " + callFnSelector(c) + "(ctx, specs)\n" +
		"return tanglHolds(b, ev, " + strconv.Quote(repr) + ")\n}"
}
