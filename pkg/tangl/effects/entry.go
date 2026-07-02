//ff:type feature=tangl type=model
//ff:what Entry — one do/undo effect summary entry reachable from an endpoint
package effects

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// Entry is one do/undo action in an endpoint's static effect closure.
type Entry struct {
	Kind string // "do" or "undo"
	Func ast.Ref
	Once bool
	Case string
	Node string
}
