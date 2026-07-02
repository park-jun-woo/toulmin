//ff:func feature=tangl type=codegen control=iteration dimension=1
//ff:what registerCaseExecs — writes each node's execution attachment in node order
package gen

import (
	"fmt"
	"strings"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// registerCaseExecs writes each node's RunOn handler (do/undo) and Run
// attachment (run), one block per node that has any execution edge, in
// case-node document order.
func registerCaseExecs(w *strings.Builder, subject string, c ast.Case, nodes map[string]nodeInfo) error {
	byNode := groupExecsByNode(c.Execs)
	for _, n := range c.Nodes {
		execs, ok := byNode[n.Name]
		if !ok {
			continue
		}
		if err := renderNodeExecBlock(w, subject, c.Name, nodes[n.Name], execs); err != nil {
			return err
		}
	}
	fmt.Fprintln(w)
	return nil
}
