//ff:func feature=tangl type=util control=sequence
//ff:what TestDecodePrefixItem — decodes a PREDICATE-DECLARE condition or a TEXT, registering Edge Ids and rejecting bad records
package stream

import (
	"testing"

	"github.com/geul-org/geul/go/pkg/toulmin"

	"github.com/geul-org/geul/go/pkg/core"
	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/packet"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// TestDecodePrefixItem covers the expression and TEXT paths, their decode
// errors, and a repeated Edge Id.
func TestDecodePrefixItem(t *testing.T) {
	d := newDec(testWire)
	ew, _ := packet.EncodeExpr(packet.ExprPacket{ID: 3, Cond: ast.Compare{Field: "level", Op: "is at least", Value: "5"}}, testWire)
	if got, err := decodePrefixItem(d, rawOf(t, ew)); err != nil || got == nil || !d.seen[3] {
		t.Fatalf("expr: %v %v", got, err)
	}
	tw, _ := toulmin.EncodeText(toulmin.TextPacket{TID: 4, Bytes: []byte("hello")}, testWire)
	if got, err := decodePrefixItem(d, rawOf(t, tw)); err != nil || got != nil || d.texts[4] != "hello" || !d.seen[4] {
		t.Fatalf("text: %v %v %v", got, err, d.texts)
	}
	if _, err := decodePrefixItem(d, rawOf(t, tw)); err == nil {
		t.Fatal("repeated Edge Id: want error")
	}
	bad := rawPacket{Class: classPayload, Rec: core.Record{Kind: core.KindEdge, Codec: toulmin.CodecPredicateDeclare}}
	if _, err := decodePrefixItem(d, bad); err == nil {
		t.Fatal("bad expression: want error")
	}
	badText := rawPacket{Class: classPayload, Rec: core.Record{Kind: core.KindEdge, Codec: core.CodecIssuer1Word, Issuer: []core.Word{0}, Format: toulmin.FormatText}}
	if _, err := decodePrefixItem(d, badText); err == nil {
		t.Fatal("bad text: want error")
	}
}
