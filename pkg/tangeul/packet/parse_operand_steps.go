//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what parseOperandSteps — read a stepped operand's [step count][label idx × step] into its step strings
package packet

import (
	"errors"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// parseOperandSteps reads a stepped operand body (a step count followed by that
// many label indices) from recs at *pos, advancing *pos and resolving each index
// through labels. It is the shared body of parseOperand for Path and SpecRef.
func parseOperandSteps(recs []word.Word, pos *int, labels []string) ([]string, error) {
	if len(recs) < *pos+1 {
		return nil, errors.New("tangeul/packet: EXPR operand step count missing")
	}
	n := int(recs[*pos])
	*pos++
	if len(recs) < *pos+n {
		return nil, errors.New("tangeul/packet: EXPR operand steps truncated")
	}
	steps := make([]string, n)
	for i := 0; i < n; i++ {
		idx := int(recs[*pos])
		if idx >= len(labels) {
			return nil, errors.New("tangeul/packet: EXPR operand step label out of range")
		}
		steps[i] = labels[idx]
		*pos++
	}
	return steps, nil
}
