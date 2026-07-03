//ff:func feature=tangl type=util control=selection
//ff:what parseCompareValue — read a Compare value slot (inline NUM32 or Text label) as its string
package packet

import (
	"errors"
	"strconv"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// parseCompareValue reads a Compare value slot from recs at *pos and returns it
// as its source string, advancing *pos. A NUM value is a 2-word NUM32 rendered
// back to decimal; a Text value is a one-word label index resolved through
// labels. An unknown ValueKind or an out-of-range index is a hard error.
func parseCompareValue(recs []word.Word, pos *int, labels []string, kind ValueKind) (string, error) {
	switch kind {
	case ValueNum:
		if len(recs) < *pos+2 {
			return "", errors.New("tangeul/packet: EXPR Compare NUM value truncated")
		}
		n, err := word.DecodeNum32(recs[*pos : *pos+2])
		if err != nil {
			return "", err
		}
		*pos += 2
		return strconv.FormatUint(uint64(n), 10), nil
	case ValueText:
		if len(recs) < *pos+1 {
			return "", errors.New("tangeul/packet: EXPR Compare Text value truncated")
		}
		idx := int(recs[*pos])
		if idx >= len(labels) {
			return "", errors.New("tangeul/packet: EXPR Compare value label out of range")
		}
		*pos++
		return labels[idx], nil
	default:
		return "", errors.New("tangeul/packet: unknown EXPR ValueKind")
	}
}
