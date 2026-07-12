//ff:func feature=tangl type=util control=sequence
//ff:what emitCompareValue — encode a literal Compare value slot as [ValueKind][inline NUM32 | Text label]
package packet

import (
	"strconv"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// emitCompareValue encodes a literal Compare value string as its value slot: an
// exact non-negative uint32 (its decimal round-trip is identical) becomes a
// ValueNum tag plus an inline 2-word NUM32; every other string becomes a
// ValueText tag plus an interned label index, so any value string round-trips
// exactly. Both the legacy kind-0 Compare record and ExprCompareX's literal
// right side reuse this, keeping their value bytes identical.
func emitCompareValue(value string, labels *[]string) []word.Word {
	if n, err := strconv.ParseUint(value, 10, 32); err == nil && strconv.FormatUint(n, 10) == value {
		num := word.EncodeNum32(uint32(n))
		return []word.Word{word.Word(ValueNum), num[0], num[1]}
	}
	return []word.Word{word.Word(ValueText), word.Word(internLabel(labels, value))}
}
