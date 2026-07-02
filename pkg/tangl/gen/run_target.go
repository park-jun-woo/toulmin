//ff:func feature=tangl type=codegen control=iteration dimension=1
//ff:what runTarget — returns the single case name a node's run edge Runs
package gen

import (
	"fmt"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// runTarget returns the single case name a node's "run" edge Runs, or ""
// if it has none. More than one "run" on the same node is rejected: the
// engine attaches at most one RunGraph per node (the spec's own
// implementation-note limitation).
func runTarget(caseName, nodeName string, execs []ast.Exec) (string, error) {
	target := ""
	for _, e := range execs {
		if e.Kind != ast.RunExec {
			continue
		}
		if target != "" {
			return "", fmt.Errorf("gen: case %q node %q has more than one `run` edge; unsupported (one RunGraph per node)", caseName, nodeName)
		}
		target = e.Case
	}
	return target, nil
}
