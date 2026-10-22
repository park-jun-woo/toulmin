//ff:func feature=tangl type=util control=sequence
//ff:what nodeHeader — pack a Node header word from Role, ObjKind and Flags for hand-built test fixtures
package packet

import (
	"github.com/geul-org/geul/go/pkg/core"
	"github.com/geul-org/geul/go/pkg/toulmin"
)

// nodeHeader packs a Node header word from Role, ObjKind and Flags, matching
// DecodeNode's bit layout (Role bits 15-14, ObjKind bits 8-7, Flags bits 6-0).
func nodeHeader(role toulmin.RoleCode, kind toulmin.ObjKind, flags toulmin.NodeFlags) core.Word {
	return core.Word(uint16(role)<<14 | uint16(kind)<<7 | uint16(flags))
}
