//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestEmitBody — encodes SOURCE plus every block and reports the last Edge Id; width overflows fail
package stream

import (
	"testing"

	"github.com/geul-org/geul/go/pkg/toulmin"
	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
)

// TestEmitBody encodes the transfer document, then forces each failure: a
// SOURCE id past the width, a block whose ids overflow, a case transition
// whose CASE-BEGIN id overflows and a closing CASE-END whose id overflows.
func TestEmitBody(t *testing.T) {
	d := buildTransferDoc()
	tt, _ := assignTIDs(d)
	ws, last, err := emitBody(d, tt, toulmin.Wire{ID: 1, Span: 1})
	if err != nil || len(ws) == 0 || last < tt.primaryEnd {
		t.Fatalf("got %d words last=%d err=%v", len(ws), last, err)
	}
	wide := *tt
	wide.primaryEnd = 0x10000
	if _, _, err := emitBody(d, &wide, toulmin.Wire{ID: 1, Span: 1}); err == nil {
		t.Fatal("SOURCE id past width: want error")
	}
	near := *tt
	near.primaryEnd = 0xFFFF
	if _, _, err := emitBody(d, &near, toulmin.Wire{ID: 1, Span: 1}); err == nil {
		t.Fatal("block ids past width: want error")
	}
	caseOnly := &doc.Doc{AST: d.AST, Source: d.Source, Blocks: nil}
	for _, b := range d.Blocks {
		if b.Kind == doc.StmtBlock && len(b.Ref.Path) == 1 && b.Ref.Sec == 5 {
			caseOnly.Blocks = append(caseOnly.Blocks, b)
		}
	}
	end := *tt
	end.primaryEnd = 0xFFFF - 3
	if _, _, err := emitBody(caseOnly, &end, toulmin.Wire{ID: 1, Span: 1}); err == nil {
		t.Fatal("case control ids past width: want error")
	}
	var lastCase []doc.Block
	for _, b := range d.Blocks {
		if b.Kind == doc.StmtBlock && b.Ref.Sec == toulmin.SectionCases && len(b.Ref.Path) == 1 {
			lastCase = []doc.Block{b}
		}
	}
	closing := *tt
	closing.primaryEnd = 0x10000 - 4 // SOURCE, CASE-BEGIN, TEXT, LABEL fit; CASE-END overflows
	if _, _, err := emitBody(&doc.Doc{AST: d.AST, Source: d.Source, Blocks: lastCase}, &closing, toulmin.Wire{ID: 1, Span: 1}); err == nil {
		t.Fatal("closing CASE-END id past width: want error")
	}
}
