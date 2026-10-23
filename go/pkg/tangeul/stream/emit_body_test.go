//ff:func feature=tangl type=util control=sequence
//ff:what TestEmitBody — encodes SOURCE plus every block; a bad SPAN width, a failing block and a recorded record error fail
package stream

import (
	"testing"

	"github.com/geul-org/geul/go/pkg/toulmin"
	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
)

// TestEmitBody encodes the transfer document, then forces each failure: an
// undeclarable SPAN width (the SOURCE fails), a block of unknown kind
// (emitBlock fails) and a PROSE span past the 1-word width (recorded on the
// Writer state and reported at the end of the walk).
func TestEmitBody(t *testing.T) {
	d := buildTransferDoc()
	tt, _ := assignTIDs(d)
	ws, err := emitBody(d, tt, testWire)
	if err != nil || len(ws) == 0 {
		t.Fatalf("got %d words err=%v", len(ws), err)
	}
	if _, err := emitBody(d, tt, toulmin.Wire{Span: 3}); err == nil {
		t.Fatal("SPAN width 3: want error")
	}
	bad := &doc.Doc{AST: d.AST, Source: d.Source, Blocks: []doc.Block{{Kind: doc.BlockKind(99)}}}
	if _, err := emitBody(bad, tt, testWire); err == nil {
		t.Fatal("unknown block kind: want error")
	}
	wide := &doc.Doc{AST: d.AST, Source: d.Source, Blocks: []doc.Block{{Kind: doc.ProseBlock, Span: doc.Span{Off: 0x10000}}}}
	if _, err := emitBody(wide, tt, testWire); err == nil {
		t.Fatal("PROSE span past the width: want error")
	}
}
