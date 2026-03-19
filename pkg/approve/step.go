//ff:type feature=approve type=model
//ff:what Step: 단일 승인 단계
package approve

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// Step holds a single approval step.
type Step struct {
	Name  string
	Graph *toulmin.Graph
}
