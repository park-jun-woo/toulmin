//ff:func feature=tangl type=util control=sequence
//ff:what execFlags — build the Edge flags and optional GATE16 from an exec's once/certainty
package stream

import (
	"github.com/geul-org/geul/go/pkg/toulmin"
	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// execFlags returns a partial EdgePacket carrying the flag word (F_ONCE) and, for
// an `if <op> <N>% certain` clause, F_GATE plus a GATE16 (§3.2). The caller fills
// in the PropCode, src, dst and span.
func execFlags(x ast.Exec) (toulmin.EdgePacket, error) {
	var ep toulmin.EdgePacket
	if x.Once {
		ep.Flags |= toulmin.EdgeFlagOnce
	}
	if x.Certainty != nil {
		op, err := codebook.GateOpFromAST(x.Certainty.Op)
		if err != nil {
			return toulmin.EdgePacket{}, err
		}
		ep.Flags |= toulmin.EdgeFlagGate
		g := toulmin.Gate{Op: op, Percent: uint8(x.Certainty.Percent)}
		ep.Gate = &g
	}
	return ep, nil
}
