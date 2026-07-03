//ff:func feature=tangl type=util control=sequence
//ff:what checkVersion — validate the STREAM_HEADER version stamp, rejecting an unknown major
package stream

import (
	"fmt"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// checkVersion decodes the STREAM_HEADER control (§7.1) and validates its
// version stamp against this build. It rejects a stream whose GEUL major, TANGL
// profile id, or codebook major differs — an unknown major means the bit layout
// or codebook may be incompatible, so decoding must fail loudly rather than
// misread. Minor/patch drift is tolerated (forward compatible).
func checkVersion(rp rawPacket) error {
	if rp.Class != word.ClassControl {
		return fmt.Errorf("tangeul/stream: stream must open with a Control STREAM_HEADER, got class %d", uint8(rp.Class))
	}
	_, head, err := word.DecodeFrame(rp.Words)
	if err != nil {
		return err
	}
	body := rp.Words[head:]
	if len(body) < 5 {
		return fmt.Errorf("tangeul/stream: STREAM_HEADER body must be at least 5 words, got %d", len(body))
	}
	if codebook.ControlCode(body[0]) != codebook.ControlStreamHeader {
		return fmt.Errorf("tangeul/stream: first control is 0x%02X, not STREAM_HEADER", uint16(body[0]))
	}
	geulMajor := uint8(uint16(body[1]) >> 8)
	profile := uint16(body[3])
	cbMajor := uint8(uint16(body[4]) >> 8)
	if geulMajor != codebook.GEULVersionMajor {
		return fmt.Errorf("tangeul/stream: unknown GEUL major version %d (this build targets %d)", geulMajor, codebook.GEULVersionMajor)
	}
	if profile != codebook.TANGLProfileID {
		return fmt.Errorf("tangeul/stream: unknown TANGL profile 0x%04X (this build is 0x%04X)", profile, codebook.TANGLProfileID)
	}
	if cbMajor != codebook.CodebookVersionMajor {
		return fmt.Errorf("tangeul/stream: unknown codebook major version %d (this build is %d)", cbMajor, codebook.CodebookVersionMajor)
	}
	return nil
}
