//ff:type feature=tangl type=model
//ff:what Registry — every loaded document in a multi-document checking graph, keyed by canonical path
package registry

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// Registry holds every ast.Document reachable from a root document via
// checking-alias See references, keyed by canonical absolute .tangeul path
// (not by alias — the same physical document may be seen under different
// aliases from different referencing documents).
type Registry struct {
	Docs map[string]*ast.Document
}
