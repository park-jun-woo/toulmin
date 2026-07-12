//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what emitOperandSteps — encode a stepped operand (Path/SpecRef) as [kind][step count][label idx × step]
package packet

import "github.com/park-jun-woo/toulmin/pkg/tangeul/word"

// emitOperandSteps encodes a stepped operand (a Path or SpecRef) as its
// OperandKind, the number of steps, then one interned label index per step. It
// is the shared body of emitOperand for the two path-shaped operand kinds.
func emitOperandSteps(kind OperandKind, steps []string, labels *[]string) []word.Word {
	out := []word.Word{word.Word(kind), word.Word(len(steps))}
	for _, s := range steps {
		out = append(out, word.Word(internLabel(labels, s)))
	}
	return out
}
