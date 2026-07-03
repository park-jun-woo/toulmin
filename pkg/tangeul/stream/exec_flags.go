//ff:func feature=tangl type=util control=sequence
//ff:what execFlags — build the Edge flags and optional GATE16 from an exec's once/certainty
package stream

import (
	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/packet"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// execFlags returns a partial EdgePacket carrying the flag word (F_ONCE) and, for
// an `if <op> <N>% certain` clause, F_GATE plus a GATE16 (§3.2). The caller fills
// in the PropCode, src, dst and span.
func execFlags(x ast.Exec) (packet.EdgePacket, error) {
	var ep packet.EdgePacket
	if x.Once {
		ep.Flags |= codebook.EdgeFlagOnce
	}
	if x.Certainty != nil {
		op, err := codebook.GateOpFromAST(x.Certainty.Op)
		if err != nil {
			return packet.EdgePacket{}, err
		}
		ep.Flags |= codebook.EdgeFlagGate
		g := word.Gate{Op: op, Percent: uint8(x.Certainty.Percent)}
		ep.Gate = &g
	}
	return ep, nil
}
