//ff:func feature=tangl type=util control=selection
//ff:what encodeClaim — encode a Claim block (subject, predicate, object, mode)
package packet

import (
	"fmt"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// encodeClaim encodes c as a Claim block in §5.2/§10 order: subject TID,
// predicate CODE, the object (1/2/3 words by c.Kind), then the mode CODE. The
// object kind itself is not emitted here — the node encoder places it in the
// Node header word (see ObjKind). An unknown Kind is a hard error.
func encodeClaim(c Claim) ([]word.Word, error) {
	out := []word.Word{word.Word(c.Subject), word.Word(c.Predicate)}
	switch c.Kind {
	case ObjTID:
		out = append(out, word.Word(c.ObjTID))
	case ObjNum:
		n := word.EncodeNum32(c.ObjNum)
		out = append(out, n[0], n[1])
	case ObjSpan:
		s := word.EncodeSpan(c.ObjSpan)
		out = append(out, s[0], s[1], s[2])
	default:
		return nil, fmt.Errorf("tangeul/packet: unknown Claim ObjKind %d", uint16(c.Kind))
	}
	return append(out, word.Word(c.Mode)), nil
}
