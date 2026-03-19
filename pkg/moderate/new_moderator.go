//ff:func feature=moderate type=engine control=sequence
//ff:what NewModerator: Moderator 생성
package moderate

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// NewModerator creates a Moderator with the given graph.
func NewModerator(g *toulmin.Graph) *Moderator {
	return &Moderator{graph: g}
}
