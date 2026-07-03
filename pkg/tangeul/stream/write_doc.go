//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what WriteDoc — serialize a Doc to a TANGEUL byte stream (header, source, ordered packets)
package stream

import (
	"errors"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// WriteDoc serializes d into a TANGEUL byte stream (§7.1): a STREAM_HEADER and
// SEGMENT_BOUNDARY, the Source Text payload, then one packet run per Block in
// stream order. Cases are wrapped in CASE_BEGIN/END via transitionCase as the
// walk enters and leaves them, and each block's packets come from emitBlock.
// Declaration order is the Blocks order, so §7.2 holds by construction.
func WriteDoc(d *doc.Doc) ([]byte, error) {
	if d == nil || d.AST == nil {
		return nil, errors.New("tangeul/stream: WriteDoc on nil Doc or AST")
	}
	tt, err := assignTIDs(d)
	if err != nil {
		return nil, err
	}
	e := &enc{tt: tt, nodeNext: tt.primaryEnd, textNext: 1}

	out := emitHeader(codebook.LocaleCodeFromTag(d.Locale))
	out = append(out, emitSource(d.Source)...)

	st := &caseState{}
	for _, b := range d.Blocks {
		trans, err := transitionCase(tt, st, b)
		if err != nil {
			return nil, err
		}
		out = append(out, trans...)
		bw, err := emitBlock(e, d, b)
		if err != nil {
			return nil, err
		}
		out = append(out, bw...)
	}
	tail, err := closeCase(st)
	if err != nil {
		return nil, err
	}
	out = append(out, tail...)
	return word.PutWords(out), nil
}
