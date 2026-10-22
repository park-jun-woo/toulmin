//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what findRel — find the first relation edge with a given predicate in a statement body
package stream

import (
	"github.com/geul-org/geul/go/pkg/toulmin"
)

// findRel returns the first relation edge in rel whose PropCode is prop, and
// whether one was found. It lets the node builders read their linking edges
// (means/has/importedFrom/sameAs/until) without looping inline.
func findRel(rel []toulmin.EdgePacket, prop toulmin.PropCode) (toulmin.EdgePacket, bool) {
	for _, e := range rel {
		if e.Prop == prop {
			return e, true
		}
	}
	return toulmin.EdgePacket{}, false
}
