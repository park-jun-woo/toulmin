//ff:type feature=tangl type=model
//ff:what ControlPacket — a decoded Control packet (SECTION / CASE_BEGIN / CASE_END / PROSE)
package packet

import (
	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// ControlPacket is a decoded Control packet (Class 11, §2.4/§7). This type
// covers the framing controls whose body is a single argument word or a SPAN:
// SECTION (Arg = SectionCode), CASE_BEGIN / CASE_END (Arg = Case TID), and PROSE
// (Span = the prose source span, §7.4). Span is non-nil exactly for PROSE; the
// stream-level controls (STREAM_HEADER / SEGMENT_BOUNDARY / CODEBOOK_DECL) carry
// version/codebook payloads and are the stream layer's concern, not this type's.
type ControlPacket struct {
	Code codebook.ControlCode
	Arg  uint16
	Span *word.Span
}
