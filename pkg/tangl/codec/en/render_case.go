//ff:func feature=tangl type=parser control=sequence
//ff:what renderCase — render a case declaration head (`in case of `X``)
package en

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// renderCase renders a tangl:Cases entry's declaration head back to
// "in case of `<name>`" (its Require/Node/Attack/Exec body items render as
// separate StmtBlocks). It is the inverse of the parser's parseCaseItem head.
func renderCase(c ast.Case) string {
	return "in case of `" + c.Name + "`"
}
