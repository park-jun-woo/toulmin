//ff:func feature=tangl type=util control=sequence
//ff:what WriteDoc — serialize a Doc as one GEUL DOC: HEAD Meta, BODY records, END
package stream

import (
	"errors"

	"github.com/geul-org/geul/go/pkg/toulmin"

	"github.com/geul-org/geul/go/pkg/core"
	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
)

// WriteDoc serializes d as exactly one GEUL t002 document (GEUL t002 TOULMIN
// 문서): DOC-*, the HEAD Meta, BODY, the SOURCE and statement records, END.
// Field widths are declared, not fixed: the SPAN width is the narrowest that
// holds the source length, and the Edge Id width the narrowest that holds the
// highest id. Ids do not depend on widths, so the body is encoded once at the
// widest id width to learn the highest id and re-encoded when a narrower width
// suffices.
func WriteDoc(d *doc.Doc) ([]byte, error) {
	if d == nil || d.AST == nil {
		return nil, errors.New("tangeul/stream: WriteDoc on nil Doc or AST")
	}
	tt, err := assignTIDs(d)
	if err != nil {
		return nil, err
	}
	w := toulmin.Wire{ID: 4, Span: core.MinWidth(uint64(len(d.Source)))}
	body, last, err := emitBody(d, tt, w)
	if err != nil {
		return nil, err
	}
	if w.ID = core.MinWidth(uint64(last)); w.ID != 4 {
		body, _, _ = emitBody(d, tt, w)
	}
	head, _ := emitHead(d.Locale, w.ID) // w.ID is always 1, 2 or 4 (MinWidth)
	end, _ := core.EncodeRecord(core.Record{Kind: core.KindMeta, Codec: core.MetaEnd})
	return core.PutWords(core.WrapDocument(append(append(head, body...), end...))), nil
}
