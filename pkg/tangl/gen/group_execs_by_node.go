//ff:func feature=tangl type=codegen control=iteration dimension=1
//ff:what groupExecsByNode — partitions a case's Execs by their trigger node
package gen

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// groupExecsByNode partitions a case's Execs by their trigger Node,
// preserving each node's original relative document order — the basis for
// "the same node's execs fire in document order" (see the spec's
// execution-order-determinism section).
func groupExecsByNode(execs []ast.Exec) map[string][]ast.Exec {
	byNode := make(map[string][]ast.Exec)
	for _, e := range execs {
		byNode[e.Node] = append(byNode[e.Node], e)
	}
	return byNode
}
