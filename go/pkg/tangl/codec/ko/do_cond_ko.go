//ff:func feature=tangl type=parser control=sequence
//ff:what doCondKo — render a do edge's condition (a bare 면 or a 확신 gate replacing it)
package ko

import "github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"

// doCondKo renders the condition that follows the node in a do edge: the bare
// 받침-selected 면/이면 when there is no confidence gate, or the 확신 gate (which
// replaces the 면) when one is present. Extracting it keeps renderExecKo's DoExec
// branch flat.
func doCondKo(e ast.Exec) (string, error) {
	if e.Certainty == nil {
		return josaFor(e.Node, "이면", "면"), nil
	}
	return renderCertaintyKo(*e.Certainty)
}
