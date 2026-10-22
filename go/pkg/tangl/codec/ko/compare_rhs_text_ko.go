//ff:func feature=tangl type=parser control=sequence
//ff:what compareRhsTextKo — render a comparison right side as a verbatim Value or an operand, with its 조사 base
package ko

import "github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"

// compareRhsTextKo renders the right side of a Compare and returns the token that
// drives 받침-based 와/과 selection: a nil Rhs is the verbatim Value string (with the
// value doubling as its own 조사 base), and a non-nil Rhs is a rendered operand. It is
// shared by renderCompareKo and renderCompareAdnKo so both keep their nesting shallow.
func compareRhsTextKo(c ast.Compare) (text string, josaBase string, err error) {
	if c.Rhs != nil {
		return renderOperandKo(c.Rhs)
	}
	return c.Value, c.Value, nil
}
