//ff:func feature=tangl type=util control=sequence
//ff:what EncodeGate — encode a Gate into a GATE16 word, rejecting percent over 100
package word

import "errors"

// EncodeGate encodes g into a single GATE16 word: op in bits 15-14, reserved
// bits 13-7 held at 0, percent in bits 6-0 (GEUL §3.2). It errors if Percent
// exceeds 100 rather than silently masking the overflow.
func EncodeGate(g Gate) (Word, error) {
	if g.Percent > gateMaxPercent {
		return 0, errors.New("tangeul/word: GATE16 percent exceeds 100")
	}
	return Word(uint16(g.Op)<<14 | uint16(g.Percent)), nil
}
