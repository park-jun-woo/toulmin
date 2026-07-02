//ff:func feature=tangl type=codegen control=iteration dimension=1
//ff:what writeDoStatements — writes one gated statement per do edge on a node
package gen

import (
	"strings"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// writeDoStatements writes one gated statement per "do" edge on this
// node, in document order, each guarded by its once/certainty clauses per
// renderDoStatement. The once-guard index counts every do on this node
// (0-based), independent of which ones actually use "once".
func writeDoStatements(w *strings.Builder, subject, caseName, nodeName string, execs []ast.Exec) {
	idx := 0
	for _, e := range execs {
		if e.Kind != ast.DoExec {
			continue
		}
		key := onceKey(subject, caseName, nodeName, idx)
		renderDoStatement(w, e, key)
		idx++
	}
}
