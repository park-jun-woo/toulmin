//ff:func feature=tangl type=validator control=selection
//ff:what certaintyGateReach — classify a do edge's certainty gate as never passing or always passing once the handler fires
package validate

import "github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"

// certaintyGateReach reports whether gate g can never pass (never) or always
// passes (always) for a handler that fires. A do handler fires only when its
// node prevails (Activated && verdict > 0), i.e. at a certainty strictly above
// 50% (% = (verdict + 1) × 50, so the reachable range is (50, 100]). Hence
// `at most N%` / `less than N%` with N ≤ 50 never fire, and `at least N%` /
// `above N%` with N ≤ 50 are redundant; thresholds past 100 are the mirror
// cases. An unknown operator is neither.
func certaintyGateReach(g ast.Certainty) (never, always bool) {
	n := g.Percent
	switch g.Op {
	case "at least":
		return n > 100, n <= 50
	case "above":
		return n >= 100, n <= 50
	case "less than":
		return n <= 50, n > 100
	case "at most":
		return n <= 50, n >= 100
	default:
		return false, false
	}
}
