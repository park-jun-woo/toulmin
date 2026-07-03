//ff:type feature=tangl type=model
//ff:what NodePacket — a decoded Node packet (frame + TID + type + header + flag-ordered slots + SPAN)
package packet

import (
	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// NodePacket is a decoded Node packet (Class 00, §5.1): the self TID, the
// NodeType, a header word packing Role (bits 15-14), the Claim object kind (bits
// 8-7) and the Flags (bits 6-0), then the conditional slots in flag order, then
// a trailing SPAN anchor. Each optional field is present exactly when its Flags
// bit is set: Func↔F_FUNC, Qual↔F_QUAL, Check↔F_CHECK, Claim↔F_CLAIM,
// Verdict↔F_VERDICT, With↔F_WITH; F_STRIKE has no slot. With holds the `with`
// clause term-name Text TIDs (§design 2). Verdict/Qual are word.VFix fixed-point
// values.
type NodePacket struct {
	TID     TID
	Type    codebook.NodeType
	Role    codebook.RoleCode
	Flags   codebook.NodeFlags
	Func    *FuncRef
	Qual    *word.VFix
	Check   *TID
	Claim   *Claim
	Verdict *word.VFix
	With    []TID
	Span    word.Span
}
