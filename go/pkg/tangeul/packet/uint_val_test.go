//ff:func feature=tangl type=util control=sequence
//ff:what uintVal — test helper: the 2-word embedded UNSIGNED-INTEGER Value of a one-word integer
package packet

import "github.com/geul-org/geul/go/pkg/core"

// uintVal returns the encoded UNSIGNED-INTEGER Value record of n (n ≤ 0xFFFF):
// the TAG WORD 0x0882 (Value, 2 words, codec 0x02) then n. It stands in for the
// old inline NUM32 in EXPR golden records.
func uintVal(n uint16) [2]core.Word {
	return [2]core.Word{0x0882, core.Word(n)}
}
