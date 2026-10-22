//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what ReadDoc — parse one GEUL DOC into a Doc: HEAD settings, SOURCE, then statement records
package stream

import (
	"github.com/geul-org/geul/go/pkg/core"
	"github.com/geul-org/geul/go/pkg/toulmin"
	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
)

// ReadDoc parses b, which must be exactly one GEUL t002 document, into a Doc —
// the inverse of WriteDoc. It validates the DOC container, reads the HEAD (spec,
// profile, Edge Id width, locale), lifts the SOURCE record, then walks the
// remaining records: framing controls (handleControl) place prose, sections and
// case membership, and each statement group rebuilds one AST element. Forward
// case-name references are patched once all names are known.
func ReadDoc(b []byte) (*doc.Doc, error) {
	ws, err := core.ReadWords(b)
	if err != nil {
		return nil, err
	}
	gd, err := core.SplitDocument(ws)
	if err != nil {
		return nil, err
	}
	h, err := readHead(gd.Head)
	if err != nil {
		return nil, err
	}
	packets, err := readPackets(gd.Body)
	if err != nil {
		return nil, err
	}
	src, spanWidth, err := readPreamble(packets, h.idWidth)
	if err != nil {
		return nil, err
	}
	d := newDec(toulmin.Wire{ID: h.idWidth, Span: spanWidth})
	if err := claimID(d, src.TID); err != nil {
		return nil, err
	}
	pendingSection := -1
	for i := 1; i < len(packets); {
		if i, err = stepRead(d, packets, i, &pendingSection); err != nil {
			return nil, err
		}
	}
	applyExecFixups(d)
	applyCheckFixups(d)
	return &doc.Doc{AST: d.ast, Source: src.Bytes, Blocks: d.blocks, Locale: h.locale}, nil
}
