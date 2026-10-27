//ff:func feature=tangl type=util control=sequence
//ff:what emitCompareValue — encode a Compare literal as an UNSIGNED-INTEGER Value or a label reference
package packet

import (
	"strconv"

	"github.com/geul-org/gogeul/pkg/core"
)

// emitCompareValue encodes a literal Compare value string as its value slot: an
// exact non-negative integer (its decimal round-trip is identical) becomes a
// ValueNum tag plus an embedded UNSIGNED-INTEGER Value; every other string
// becomes a ValueText tag plus an interned label index, so any value string
// round-trips exactly. The Compare and ExprCompareX records both reuse this.
func emitCompareValue(value string, labels *[]string) []core.Word {
	if n, err := strconv.ParseUint(value, 10, 64); err == nil && strconv.FormatUint(n, 10) == value {
		v, _ := core.EncodeRecord(core.EncodeUintValue(n))
		return append([]core.Word{core.Word(ValueNum)}, v...)
	}
	return []core.Word{core.Word(ValueText), core.Word(internLabel(labels, value))}
}
