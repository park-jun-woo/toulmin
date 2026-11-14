//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what textsOf — the strings of a list of TEXT Edge Ids, in order
package stream

import "github.com/geul-org/gogeul/pkg/toulmin"

// textsOf resolves TEXT Edge Ids (a node's F_WITH or F_GROUNDS list) to their
// strings in order; nil for an empty list.
func textsOf(d *dec, ids []toulmin.TID) []string {
	var out []string
	for _, id := range ids {
		out = append(out, d.texts[id])
	}
	return out
}
