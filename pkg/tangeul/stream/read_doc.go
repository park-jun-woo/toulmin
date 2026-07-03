//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what ReadDoc — parse a TANGEUL byte stream back into a Doc (AST + Source + Blocks)
package stream

import (
	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// ReadDoc parses a TANGEUL byte stream into a Doc, the inverse of WriteDoc. It
// checks the STREAM_HEADER version stamp, lifts the Source Text payload, then
// walks the remaining packets: controls (handled by handleControl) frame prose,
// section and case membership, and each statement group rebuilds one AST element.
// Forward case-name references are patched once all names are known.
func ReadDoc(b []byte) (*doc.Doc, error) {
	ws, err := word.ReadWords(b)
	if err != nil {
		return nil, err
	}
	packets, err := readPackets(ws)
	if err != nil {
		return nil, err
	}
	src, localeCode, err := readPreamble(packets)
	if err != nil {
		return nil, err
	}

	d := newDec()
	pendingSection := -1
	for i := 3; i < len(packets); {
		next, err := stepRead(d, packets, i, &pendingSection)
		if err != nil {
			return nil, err
		}
		i = next
	}
	applyExecFixups(d)
	applyCheckFixups(d)
	return &doc.Doc{AST: d.ast, Source: src.Bytes, Blocks: d.blocks, Locale: codebook.LocaleTagFromCode(localeCode)}, nil
}
