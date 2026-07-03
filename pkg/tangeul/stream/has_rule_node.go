//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what hasRuleNode — report whether a statement's primary nodes include a rule node
package stream

import (
	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/packet"
)

// hasRuleNode reports whether primaries contains a Rule node. Combined with the
// open-case state, it tells the reader whether a statement's leading payload is a
// Rules condition EXPR.
func hasRuleNode(primaries []packet.NodePacket) bool {
	for _, p := range primaries {
		if p.Type == codebook.NodeRule {
			return true
		}
	}
	return false
}
