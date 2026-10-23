//ff:func feature=tangl type=util control=sequence
//ff:what proseWords — encode a PROSE record for a source span
package stream

import (
	"github.com/geul-org/geul/go/pkg/core"
	"github.com/geul-org/geul/go/pkg/toulmin"
	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
)

// proseWords encodes a PROSE record carrying s. It backs both standalone prose
// blocks and the heading span that follows a SECTION.
func proseWords(e *enc, s doc.Span) []core.Word {
	ws := toWireSpan(s)
	return controlWords(e, toulmin.ControlProse, 0, &ws)
}
