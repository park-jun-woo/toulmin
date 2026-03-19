//ff:type feature=moderate type=engine
//ff:what Moderator: 모더레이션 판정 실행
package moderate

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// Moderator evaluates content moderation rules.
type Moderator struct {
	graph *toulmin.Graph
}
