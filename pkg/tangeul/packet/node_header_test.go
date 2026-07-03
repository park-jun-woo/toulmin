//ff:func feature=tangl type=util control=sequence
//ff:what nodeHeader — pack a Node header word from Role, ObjKind and Flags for hand-built test fixtures
package packet

import (
	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// nodeHeader packs a Node header word from Role, ObjKind and Flags, matching
// DecodeNode's bit layout (Role bits 15-14, ObjKind bits 8-7, Flags bits 6-0).
func nodeHeader(role codebook.RoleCode, kind ObjKind, flags codebook.NodeFlags) word.Word {
	return word.Word(uint16(role)<<14 | uint16(kind)<<7 | uint16(flags))
}
