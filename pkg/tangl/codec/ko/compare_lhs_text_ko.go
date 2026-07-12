//ff:func feature=tangl type=parser control=sequence
//ff:what compareLhsTextKo — render a comparison left side as bare `field` or an operand, with its 조사 base
package ko

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// compareLhsTextKo renders the left side of a Compare and returns the token that
// drives 받침-based 조사 selection: a nil Lhs is the bare backtick `field` (the
// byte-identical path for pre-existing documents), and a non-nil Lhs is a rendered
// operand. It is shared by renderCompareKo and renderCompareAdnKo so both keep their
// nesting shallow.
func compareLhsTextKo(c ast.Compare) (text string, josaBase string, err error) {
	if c.Lhs != nil {
		return renderOperandKo(c.Lhs)
	}
	return "`" + c.Field + "`", c.Field, nil
}
