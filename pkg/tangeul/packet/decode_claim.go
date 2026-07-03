//ff:func feature=tangl type=util control=selection
//ff:what decodeClaim — decode a Claim block given its object kind, returning words consumed
package packet

import (
	"errors"
	"fmt"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// decodeClaim decodes a Claim block from ws given kind (read by the caller from
// the Node header word). It returns the Claim and the number of words consumed
// (4, 5 or 6 by kind). A short slice or an unknown kind is a hard error.
func decodeClaim(ws []word.Word, kind ObjKind) (Claim, int, error) {
	if len(ws) < 2 {
		return Claim{}, 0, errors.New("tangeul/packet: Claim requires at least a subject and predicate")
	}
	c := Claim{Subject: TID(ws[0]), Predicate: uint16(ws[1]), Kind: kind}
	pos := 2
	switch kind {
	case ObjTID:
		if len(ws) < pos+1 {
			return Claim{}, 0, errors.New("tangeul/packet: Claim TID object requires 1 word")
		}
		c.ObjTID = TID(ws[pos])
		pos++
	case ObjNum:
		n, err := word.DecodeNum32(ws[pos:])
		if err != nil {
			return Claim{}, 0, err
		}
		c.ObjNum = n
		pos += 2
	case ObjSpan:
		s, err := word.DecodeSpan(ws[pos:])
		if err != nil {
			return Claim{}, 0, err
		}
		c.ObjSpan = s
		pos += 3
	default:
		return Claim{}, 0, fmt.Errorf("tangeul/packet: unknown Claim ObjKind %d", uint16(kind))
	}
	if len(ws) < pos+1 {
		return Claim{}, 0, errors.New("tangeul/packet: Claim requires a mode word")
	}
	c.Mode = uint16(ws[pos])
	return c, pos + 1, nil
}
