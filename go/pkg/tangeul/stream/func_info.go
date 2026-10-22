//ff:type feature=tangl type=model
//ff:what funcInfo — a decoded Func node's package alias TID and its function name
package stream

import (
	"github.com/geul-org/geul/go/pkg/toulmin"
)

// funcInfo is what the reader records for each materialized Func node (§2.1):
// the package alias TID its F_FUNC slot bound (0 for a local reference) and the
// function name recovered from its label. A referring node/edge whose FuncRef
// Alias points at the Func node resolves back to an ast.Ref through this.
type funcInfo struct {
	alias toulmin.TID
	name  string
}
