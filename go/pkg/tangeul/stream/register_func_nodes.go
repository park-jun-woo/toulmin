//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what registerFuncNodes — record each Func node's package alias and name for later resolution
package stream

import (
	"github.com/geul-org/gogeul/pkg/toulmin"
)

// registerFuncNodes records every materialized Func node of a statement into the
// funcRef map: its package alias TID (from the F_FUNC slot) and its name (its
// label), so a referring node/edge can rebuild the ast.Ref.
func registerFuncNodes(d *dec, funcNodes []toulmin.NodePacket, labels map[toulmin.TID][]string) {
	for _, fn := range funcNodes {
		var alias toulmin.TID
		if fn.Func != nil {
			alias = toulmin.TID(fn.Func.Alias)
		}
		d.funcRef[toulmin.TID(fn.TID)] = funcInfo{alias: alias, name: firstLabel(labels, toulmin.TID(fn.TID))}
	}
}
