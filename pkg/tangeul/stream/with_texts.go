//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what withTexts — issue a Text payload per `with` term, returning their TIDs for the F_WITH slot
package stream

import (
	"github.com/park-jun-woo/toulmin/pkg/tangeul/packet"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// withTexts issues one Text payload per `with` term name onto the prefix and
// returns their TIDs, which fill a rule node's F_WITH slot (§5.5). The terms are
// Definitions names carried as strings, so no term-node linkage is needed.
func withTexts(e *enc, pre *[]word.Word, terms []string) []packet.TID {
	out := make([]packet.TID, 0, len(terms))
	for _, t := range terms {
		out = append(out, packet.TID(addText(e, pre, t)))
	}
	return out
}
