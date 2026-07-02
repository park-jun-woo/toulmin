//ff:func feature=tangl type=codegen control=iteration dimension=1
//ff:what writeNodeExecs — writes a node's do/undo edges in document order, arming each undo right after its preceding dos succeed
package gen

import (
	"strings"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// writeNodeExecs writes this node's "do" and "undo" edges in document
// order. Each "do" is a gated statement (per renderDoStatement); each
// "undo" is pushed onto the compensation stack immediately at that point
// in the sequence. This keeps a do/undo/do interleaving on the same node
// faithful to the spec's undo-arming rule — compensation for an earlier
// do is armed before a later do on the same node can fail — rather than
// arming every undo only after all of the node's dos have run. The
// once-guard index counts every do on this node (0-based), independent of
// which ones actually use "once".
func writeNodeExecs(w *strings.Builder, subject, caseName, nodeName string, execs []ast.Exec) {
	doIdx := 0
	for _, e := range execs {
		switch e.Kind {
		case ast.DoExec:
			key := onceKey(subject, caseName, nodeName, doIdx)
			renderDoStatement(w, e, key)
			doIdx++
		case ast.UndoExec:
			writeUndoPush(w, e)
		}
	}
}
