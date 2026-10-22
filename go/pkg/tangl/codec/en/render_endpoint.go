//ff:func feature=tangl type=parser control=sequence
//ff:what renderEndpoint — render an endpoint declaration head (`provides `X``)
package en

import "github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"

// renderEndpoint renders a tangl:Provides entry's declaration head back to
// "provides `<name>`" (its Require/Run/Check body items render as separate
// StmtBlocks). It is the inverse of the parser's parseEndpointItem head.
func renderEndpoint(e ast.Endpoint) string {
	return "provides `" + e.Name + "`"
}
