//ff:func feature=tangl type=util control=sequence
//ff:what decodePrefixItem — decode one prefix payload as a Rules EXPR or a Text (assigning its TID)
package stream

import (
	"github.com/park-jun-woo/toulmin/pkg/tangeul/packet"
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// decodePrefixItem decodes one prefix payload: as an EXPR (returning its
// condition) when isExpr is set, otherwise as a Text whose bytes are stored under
// the next issued Text TID. Text TIDs advance in stream order, matching the
// Writer so references resolve.
func decodePrefixItem(d *dec, rp rawPacket, isExpr bool) (ast.Expr, error) {
	if isExpr {
		ep, _, err := packet.DecodeExpr(rp.Words)
		if err != nil {
			return nil, err
		}
		return ep.Cond, nil
	}
	tp, _, err := packet.DecodeText(rp.Words)
	if err != nil {
		return nil, err
	}
	d.texts[d.textNext] = string(tp.Bytes)
	d.textNext++
	return nil, nil
}
