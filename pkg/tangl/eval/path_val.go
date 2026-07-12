//ff:func feature=tangl type=engine control=sequence
//ff:what pathVal — resolves a Path operand to (value, ok) via binders or ctx
package eval

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// pathVal resolves a Path to a value. Its first step is the bound element when
// it names an in-scope exists binder, else a ctx key; remaining steps walk into
// the value via map-key or struct-field access. A missing key/field or nil
// intermediate yields ok=false.
func pathVal(p ast.Path, sc *scope) (any, bool) {
	if len(p.Steps) == 0 {
		return nil, false
	}
	first := p.Steps[0]
	var cur any
	if b, ok := sc.binders[first]; ok {
		cur = b.value
	} else {
		v, ok := sc.ctx.Get(first)
		if !ok {
			return nil, false
		}
		cur = v
	}
	return walkSteps(cur, p.Steps[1:])
}
