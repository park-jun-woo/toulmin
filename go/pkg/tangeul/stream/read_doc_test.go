//ff:func feature=tangl type=util control=sequence
//ff:what TestReadDoc — ReadDoc reads a SOURCE-only document and fails on a bad record mid-BODY or a repeated SOURCE id
package stream

import (
	"testing"

	"github.com/geul-org/gogeul/pkg/toulmin"

	"github.com/geul-org/gogeul/pkg/core"
)

// TestReadDoc builds documents by hand from emitHead: SOURCE alone (no blocks),
// SOURCE followed by a control whose payload is short (a stepRead error), and a
// SOURCE with Edge Id 0 (rejected by claimID).
func TestReadDoc(t *testing.T) {
	build := func(src toulmin.TID, extra ...core.Word) []byte {
		head := emitHead("en")
		sw, _ := toulmin.EncodeSource(toulmin.TextPacket{TID: src, Bytes: []byte("x")}, testWire)
		ws := core.WrapDocument(append(append(append(head, sw...), extra...), 0x180F))
		return core.PutWords(ws)
	}
	got, err := ReadDoc(build(1))
	if err != nil || got.AST == nil || len(got.Blocks) != 0 || got.Locale != "en" || string(got.Source) != "x" {
		t.Fatalf("SOURCE only: %+v %v", got, err)
	}
	short, _ := core.EncodeRecord(core.Record{Kind: core.KindEdge, Codec: core.CodecIssuer1Word, Issuer: []core.Word{0}, Format: toulmin.FormatCaseEnd, Payload: []core.Word{2}})
	if _, err := ReadDoc(build(1, short...)); err == nil {
		t.Fatal("bad control mid-BODY: want error")
	}
	if _, err := ReadDoc(build(0)); err == nil {
		t.Fatal("SOURCE Edge Id 0: want error")
	}

	head := emitHead("")
	tw, _ := toulmin.EncodeText(toulmin.TextPacket{TID: 1}, testWire)
	noSource := core.WrapDocument(append(append(head, tw...), 0x180F))
	if _, err := ReadDoc(core.PutWords(noSource)); err == nil {
		t.Fatal("BODY not opened by SOURCE: want error")
	}
}
