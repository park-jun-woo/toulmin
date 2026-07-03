//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what findRel — find the first relation edge with a given predicate in a statement body
package stream

import (
	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/packet"
)

// findRel returns the first relation edge in rel whose PropCode is prop, and
// whether one was found. It lets the node builders read their linking edges
// (means/has/importedFrom/sameAs/until) without looping inline.
func findRel(rel []packet.EdgePacket, prop codebook.PropCode) (packet.EdgePacket, bool) {
	for _, e := range rel {
		if e.Prop == prop {
			return e, true
		}
	}
	return packet.EdgePacket{}, false
}
