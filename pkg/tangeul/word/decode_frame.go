//ff:func feature=tangl type=util control=sequence
//ff:what DecodeFrame — decode a frame header, reporting prefix violations and consumed word count
package word

import "errors"

// DecodeFrame decodes a frame header from the head of ws. It returns the Frame,
// the number of words consumed (1, or 3 when the 32-bit extension is present),
// and an error. The Free prefix (0001) is enforced; any other prefix is a hard
// error rather than a silent acceptance. When Len equals the 1023 sentinel the
// real body word count is read from the trailing NUM32.
func DecodeFrame(ws []Word) (Frame, int, error) {
	if len(ws) < 1 {
		return Frame{}, 0, errors.New("tangeul/word: frame header requires at least 1 word")
	}
	head := ws[0]
	if head>>12 != framePrefix {
		return Frame{}, 0, errors.New("tangeul/word: frame prefix mismatch (expected 0001 Free)")
	}
	class := Class((head >> 10) & 0x3)
	len10 := uint32(head & 0x3FF)
	if len10 != frameLenSentinel {
		return Frame{Class: class, Len: len10}, 1, nil
	}
	if len(ws) < 3 {
		return Frame{}, 0, errors.New("tangeul/word: 32-bit length extension requires a trailing NUM32")
	}
	n, err := DecodeNum32(ws[1:3])
	if err != nil {
		return Frame{}, 0, err
	}
	return Frame{Class: class, Len: n}, 3, nil
}
