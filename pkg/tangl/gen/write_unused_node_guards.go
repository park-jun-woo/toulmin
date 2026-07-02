//ff:func feature=tangl type=codegen control=iteration dimension=1
//ff:what writeUnusedNodeGuards — discards node vars no Attacks/Exec edge otherwise references
package gen

import (
	"fmt"
	"strings"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// writeUnusedNodeGuards writes "_ = <var>" for every node whose only use
// is its own registration statement (no don't/do/run edge attaches to
// it) — a bare `checking` or judgment-only node, for instance — so the Go
// compiler never rejects it as an unused local variable.
func writeUnusedNodeGuards(w *strings.Builder, c ast.Case, nodes map[string]nodeInfo) {
	used := usedNodeNames(c)
	for _, n := range c.Nodes {
		if used[n.Name] {
			continue
		}
		fmt.Fprintf(w, "\t_ = %s\n", nodes[n.Name].Var)
	}
}
