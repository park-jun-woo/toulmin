//ff:func feature=tangl type=util control=sequence
//ff:what decodePrefixItem — decode one statement-prefix record: a PREDICATE-DECLARE condition or a TEXT
package stream

import (
	"github.com/geul-org/gogeul/pkg/toulmin"
	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/packet"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// decodePrefixItem decodes one prefix record by its codec: a PREDICATE-DECLARE
// is the Rules condition and its tree is returned; anything else is a TEXT
// whose string is recorded under its Edge Id. Either way its Edge Id is
// registered.
func decodePrefixItem(d *dec, rp rawPacket) (ast.Expr, error) {
	if rp.Rec.Codec == toulmin.CodecPredicateDeclare {
		ep, err := packet.DecodeExpr(rp.Rec, d.wire)
		if err != nil {
			return nil, err
		}
		return ep.Cond, claimID(d, ep.ID)
	}
	tp, err := toulmin.DecodeText(rp.Rec, d.wire)
	if err != nil {
		return nil, err
	}
	d.texts[tp.TID] = string(tp.Bytes)
	return nil, claimID(d, tp.TID)
}
