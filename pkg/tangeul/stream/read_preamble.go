//ff:func feature=tangl type=util control=sequence
//ff:what readPreamble — validate the header and boundary and return the Source Text payload and authoring locale
package stream

import (
	"errors"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/packet"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// readPreamble validates the three leading packets — STREAM_HEADER (version
// stamp + authoring locale), SEGMENT_BOUNDARY, and the Source Text payload — and
// returns the decoded Source together with the authoring locale code (§11-11). A
// stream with fewer than three packets, a bad version, or a missing source is a
// hard error.
func readPreamble(packets []rawPacket) (packet.TextPacket, uint16, error) {
	if len(packets) < 3 {
		return packet.TextPacket{}, codebook.LocaleUnknown, errors.New("tangeul/stream: stream too short (missing header, boundary, or source)")
	}
	if err := checkVersion(packets[0]); err != nil {
		return packet.TextPacket{}, codebook.LocaleUnknown, err
	}
	locale, err := readHeaderLocale(packets[0])
	if err != nil {
		return packet.TextPacket{}, codebook.LocaleUnknown, err
	}
	if err := checkSegmentBoundary(packets[1]); err != nil {
		return packet.TextPacket{}, codebook.LocaleUnknown, err
	}
	if packets[2].Class != word.ClassPayload {
		return packet.TextPacket{}, codebook.LocaleUnknown, errors.New("tangeul/stream: third packet must be the Source Text payload")
	}
	src, _, err := packet.DecodeText(packets[2].Words)
	return src, locale, err
}
