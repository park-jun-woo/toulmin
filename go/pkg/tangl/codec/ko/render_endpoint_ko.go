//ff:func feature=tangl type=parser control=sequence
//ff:what renderEndpointKo — render an endpoint declaration head (`` `X` 를 한다면 ``)
package ko

import "github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"

// renderEndpointKo renders a tangl:업무 entry's declaration head back to "`<name>` 를 한다면"
// (its 필수다/진행한다/판정한다 body items render as separate StmtBlocks). It is the inverse
// of parseEndpointKo's head; the 를(을/를) particle is 받침-selected and 한다면 is one
// keyword (its 면 is not a standalone particle).
func renderEndpointKo(e ast.Endpoint) string {
	return "`" + e.Name + "` " + josaFor(e.Name, "을", "를") + " 한다면"
}
