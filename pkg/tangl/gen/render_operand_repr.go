//ff:func feature=tangl type=codegen control=selection
//ff:what renderOperandRepr — renders an Operand's Where notation as a reprCode
package gen

import (
	"strconv"
	"strings"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// renderOperandRepr renders an Operand's Where notation (the string
// pkg/tangl/eval.reprOperand produces at runtime) as a reprCode: a Path or
// SpecRef becomes its dotted-steps constant, a Call its "alias.name" or bare
// "name" constant, and a binder-headed Path a runtime concatenation
// "<source>[<idx>].<rest>" indexed by the enclosing element closure's index
// variable.
func renderOperandRepr(o ast.Operand, env map[string]binderInfo) reprCode {
	switch v := o.(type) {
	case ast.Path:
		if len(v.Steps) == 0 {
			return reprCode{Expr: `""`, Const: true, Val: ""}
		}
		b, bound := env[v.Steps[0]]
		if !bound {
			val := strings.Join(v.Steps, ".")
			return reprCode{Expr: strconv.Quote(val), Const: true, Val: val}
		}
		rest := ""
		if len(v.Steps) > 1 {
			rest = "." + strings.Join(v.Steps[1:], ".")
		}
		if b.Src.Const {
			expr := strconv.Quote(b.Src.Val+"[") + "+strconv.Itoa(" + b.Idx + ")+" + strconv.Quote("]"+rest)
			return reprCode{Expr: expr}
		}
		expr := b.Src.Expr + `+"["+strconv.Itoa(` + b.Idx + `)+`
		expr += strconv.Quote("]" + rest)
		return reprCode{Expr: expr}
	case ast.SpecRef:
		val := strings.Join(v.Steps, ".")
		return reprCode{Expr: strconv.Quote(val), Const: true, Val: val}
	case ast.Call:
		val := v.Fn.Name
		if v.Fn.Alias != "" {
			val = v.Fn.Alias + "." + v.Fn.Name
		}
		return reprCode{Expr: strconv.Quote(val), Const: true, Val: val}
	default:
		return reprCode{Expr: `""`, Const: true, Val: ""}
	}
}
