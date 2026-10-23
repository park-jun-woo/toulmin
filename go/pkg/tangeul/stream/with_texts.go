//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what withTexts — issue a Text payload per `with` term, returning their TIDs for the F_WITH slot
package stream

import (
	"github.com/geul-org/geul/go/pkg/core"
	"github.com/geul-org/geul/go/pkg/toulmin"
)

// withTexts issues one TEXT record per `with` term name onto the prefix and
// returns their Edge Ids, which fill a rule node's F_WITH slot (§5.5). The terms are
// Definitions names carried as strings, so no term-node linkage is needed.
func withTexts(e *enc, pre *[]core.Word, terms []string) []toulmin.TID {
	out := make([]toulmin.TID, 0, len(terms))
	for _, t := range terms {
		out = append(out, addText(e, pre, t))
	}
	return out
}
