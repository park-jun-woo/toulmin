//ff:func feature=tangl type=util control=sequence
//ff:what resolveFuncRef — rebuild an ast.Ref from a Func node's recorded alias TID and name
package stream

import (
	"fmt"

	"github.com/geul-org/geul/go/pkg/toulmin"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// resolveFuncRef rebuilds the ast.Ref a referring node/edge pointed at through a
// Func node TID: the recorded package alias TID resolves to its alias name
// (empty for a local reference), and the function name is the Func node's label.
// An unregistered Func node TID is a hard error — the stream is malformed.
func resolveFuncRef(d *dec, fTID toulmin.TID) (ast.Ref, error) {
	fi, ok := d.funcRef[fTID]
	if !ok {
		return ast.Ref{}, fmt.Errorf("tangeul/stream: function reference to unknown Func node %d", fTID)
	}
	var alias string
	if fi.alias != 0 {
		alias = d.nodeName[fi.alias]
	}
	return ast.Ref{Alias: alias, Name: fi.name}, nil
}
