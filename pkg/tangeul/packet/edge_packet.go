//ff:type feature=tangl type=model
//ff:what EdgePacket — a decoded Edge packet (frame + PropCode + src + dst + flags + GATE16 + SPAN)
package packet

import (
	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// EdgePacket is a decoded Edge packet (Class 01 = triple, §6.1): the PropCode
// predicate, the src TID, the dst, a flag word (4 flags in bits 3-0), an
// optional GATE16 (present iff F_GATE), and a trailing SPAN. The dst slot kind
// is decided by the PropCode alone via codebook.PropDstKind: a FuncRef dst
// (do/undo) uses DstFunc (2 words); every TID or Text dst uses DstTID (1 word).
// Provenance (DstBlock) is not encoded in v0.1 (README background).
type EdgePacket struct {
	Prop    codebook.PropCode
	Src     TID
	DstTID  TID
	DstFunc *FuncRef
	Flags   codebook.EdgeFlags
	Gate    *word.Gate
	Span    word.Span
}
