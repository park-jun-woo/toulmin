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
// The SPAN width is the narrowest that holds the source length; Edge Ids and
// REFs are LOCAL Reference records in minimal words, so no id width is declared.
func WriteDoc(d *doc.Doc) ([]byte, error) {
	if d == nil || d.AST == nil {
		return nil, errors.New("tangeul/stream: WriteDoc on nil Doc or AST")
	}
	tt, err := assignTIDs(d)
	if err != nil {
		return nil, err
	}
	body, err := emitBody(d, tt, toulmin.Wire{Span: core.MinWidth(uint64(len(d.Source)))})
	if err != nil {
		return nil, err
	}
	head := emitHead(d.Locale)
	end, _ := core.EncodeRecord(core.Record{Kind: core.KindMeta, Codec: core.MetaEnd})
	return core.PutWords(core.WrapDocument(append(append(head, body...), end...))), nil
}
