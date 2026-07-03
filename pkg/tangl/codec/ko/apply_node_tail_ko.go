//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what applyNodeTailKo — parse a node's trailing '기준'(with) and '가중'(qualified) clauses
package ko

import (
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// applyNodeTailKo parses the optional trailing clauses after 규칙이다: '기준 <terms>'
// (with clause) and '가중 <float>' (qualified). The clauses are comma-separated;
// skipParticlesKo peels the commas and interposed particles, and applyNodeClauseKo
// consumes one clause per iteration.
func applyNodeTailKo(n *ast.Node, tail string, line int, path string) error {
	t := skipParticlesKo(tail)
	for t != "" {
		rest, err := applyNodeClauseKo(n, t, line, path)
		if err != nil {
			return err
		}
		t = skipParticlesKo(rest)
	}
	return nil
}
