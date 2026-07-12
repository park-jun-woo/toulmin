//ff:func feature=tangl type=util control=selection
//ff:what emitOperand — encode an ast.Operand (Path/SpecRef/Call) into EXPR record words, interning labels
package packet

import (
	"fmt"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// emitOperand encodes op into a self-contained operand sub-record (§design 2),
// appending any step/alias/name strings into *labels. A Path or SpecRef writes
// its OperandKind, a step count, then one interned label index per step; a Call
// writes OperandKind, the interned alias index (the empty string for a local
// call) and the interned name index. parseOperand is its exact inverse.
func emitOperand(op ast.Operand, labels *[]string) ([]word.Word, error) {
	switch v := op.(type) {
	case ast.Path:
		return emitOperandSteps(OperandPath, v.Steps, labels), nil
	case ast.SpecRef:
		return emitOperandSteps(OperandSpec, v.Steps, labels), nil
	case ast.Call:
		return []word.Word{
			word.Word(OperandCall),
			word.Word(internLabel(labels, v.Fn.Alias)),
			word.Word(internLabel(labels, v.Fn.Name)),
		}, nil
	default:
		return nil, fmt.Errorf("tangeul/packet: unknown ast.Operand type %T", op)
	}
}
