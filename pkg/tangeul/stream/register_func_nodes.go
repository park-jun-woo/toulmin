//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what registerFuncNodes — record each Func node's package alias and name for later resolution
package stream

import "github.com/park-jun-woo/toulmin/pkg/tangeul/packet"

// registerFuncNodes records every materialized Func node of a statement into the
// funcRef map: its package alias TID (from the F_FUNC slot) and its name (its
// label), so a referring node/edge can rebuild the ast.Ref.
func registerFuncNodes(d *dec, funcNodes []packet.NodePacket, labels map[uint16][]string) {
	for _, fn := range funcNodes {
		var alias uint16
		if fn.Func != nil {
			alias = uint16(fn.Func.Alias)
		}
		d.funcRef[uint16(fn.TID)] = funcInfo{alias: alias, name: firstLabel(labels, uint16(fn.TID))}
	}
}
