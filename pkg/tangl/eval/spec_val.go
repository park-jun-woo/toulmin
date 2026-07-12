//ff:func feature=tangl type=engine control=sequence
//ff:what specVal — resolves a SpecRef operand to (value, ok) via Specs.Find
package eval

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// specVal resolves a SpecRef to a value. Its first step names a spec term
// (matched via Specs.Find); with no remaining steps the value is the Spec
// itself, otherwise the remaining steps walk into it like a Path.
func specVal(sr ast.SpecRef, sc *scope) (any, bool) {
	if len(sr.Steps) == 0 {
		return nil, false
	}
	sp := sc.specs.Find(sr.Steps[0])
	if sp == nil {
		return nil, false
	}
	if len(sr.Steps) == 1 {
		return sp, true
	}
	return walkSteps(sp, sr.Steps[1:])
}
