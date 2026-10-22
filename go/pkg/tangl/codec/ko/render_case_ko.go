//ff:func feature=tangl type=parser control=sequence
//ff:what renderCaseKo — render a case declaration head (`` `X` 의 경우에 ``)
package ko

import "github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"

// renderCaseKo renders a tangl:사안 entry's declaration head back to "`<name>` 의 경우에"
// (its 필수다/규칙이다/무효다/실행한다 body items render as separate StmtBlocks). It is the
// inverse of parseCaseKo's head; 경우에 is one keyword (its 에 is not a standalone particle).
func renderCaseKo(c ast.Case) string {
	return "`" + c.Name + "` 의 경우에"
}
