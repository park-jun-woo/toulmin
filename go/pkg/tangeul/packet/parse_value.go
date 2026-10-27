//ff:func feature=tangl type=util control=selection
//ff:what parseCompareValue — read a Compare value slot (UNSIGNED-INTEGER Value or label index) as its string
package packet

import (
	"errors"
	"strconv"

	"github.com/geul-org/gogeul/pkg/toulmin"

	"github.com/geul-org/gogeul/pkg/core"
)

// parseCompareValue reads a Compare value slot from recs at *pos and returns it
// as its source string, advancing *pos. A ValueNum is an embedded
// UNSIGNED-INTEGER Value rendered back to decimal; a ValueText is a one-word
// label index resolved through labels. An unknown ValueKind or an out-of-range
// index is an error.
func parseCompareValue(recs []core.Word, pos *int, labels []string, kind ValueKind) (string, error) {
	switch kind {
	case ValueNum:
		r, err := toulmin.TakeValue(recs, pos)
		if err != nil {
			return "", err
		}
		n, err := core.DecodeUintValue(r)
		return strconv.FormatUint(n, 10), err
	case ValueText:
		idx, err := toulmin.TakeWord(recs, pos, "EXPR Compare Text value")
		if err != nil {
			return "", err
		}
		if int(idx) >= len(labels) {
			return "", errors.New("tangeul/packet: EXPR Compare value label out of range")
		}
		return labels[idx], nil
	default:
		return "", errors.New("tangeul/packet: unknown EXPR ValueKind")
	}
}
