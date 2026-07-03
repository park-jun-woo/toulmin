//ff:func feature=tangl type=util control=sequence
//ff:what DecodeGate — decode a GATE16 word, enforcing zero reserved bits and percent range
package word

import "errors"

// DecodeGate decodes a GATE16 word into a Gate (GEUL §3.2). It enforces that
// the seven reserved bits (13-7) are zero and that the percent (bits 6-0) does
// not exceed 100; either violation is a hard error, never a silent acceptance.
func DecodeGate(w Word) (Gate, error) {
	if (w>>7)&0x7F != 0 {
		return Gate{}, errors.New("tangeul/word: GATE16 reserved bits must be 0")
	}
	percent := uint8(w & 0x7F)
	if percent > gateMaxPercent {
		return Gate{}, errors.New("tangeul/word: GATE16 percent exceeds 100")
	}
	return Gate{
		Op:      GateOp(w >> 14),
		Percent: percent,
	}, nil
}
