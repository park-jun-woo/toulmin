//ff:func feature=tangl type=parser control=sequence
//ff:what parseCompareLhsKo — read a comparison left side as an operand, demoting a single-step path to Field
package ko

import (
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// parseCompareLhsKo reads the left side of a KO comparison (the text before the
// operator suffix has been peeled). It returns a seed ast.Compare with either Field
// set (a single-step backtick path is demoted to Field, Lhs nil, so pre-existing
// documents keep a byte-identical AST) or Lhs set (a multi-step path, spec ref, or
// host call), plus the text after the operand. ok is false when no operand introducer
// was found (the caller reports the "expected backtick-quoted field" error with the
// full condition); a committed-but-malformed operand surfaces as err. It is shared by
// parseCompareKo (서술형) and parseCompareAdnKo (관형형).
func parseCompareLhsKo(before string, line int, path string) (cmp ast.Compare, rest string, ok bool, err error) {
	lhs, rest, ok, err := parseOperandKo(before, line, path)
	if err != nil || !ok {
		return ast.Compare{}, before, ok, err
	}
	if p, isPath := lhs.(ast.Path); isPath && len(p.Steps) == 1 {
		cmp.Field = p.Steps[0]
	} else {
		cmp.Lhs = lhs
	}
	return cmp, rest, true, nil
}
