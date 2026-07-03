//ff:func feature=tangl type=util control=sequence
//ff:what readHeaderLocale — extract the authoring locale code from the STREAM_HEADER packet
package stream

import (
	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// readHeaderLocale returns the authoring locale code stamped in the STREAM_HEADER
// (§11-11). The header is already shape-validated by checkVersion; here the frame
// is re-decoded only to read the 6th body word. A pre-locale header (5-word body)
// has no code, so the locale degrades to LocaleUnknown rather than erroring.
func readHeaderLocale(rp rawPacket) (uint16, error) {
	_, head, err := word.DecodeFrame(rp.Words)
	if err != nil {
		return codebook.LocaleUnknown, err
	}
	body := rp.Words[head:]
	if len(body) < 6 {
		return codebook.LocaleUnknown, nil
	}
	return uint16(body[5]), nil
}
