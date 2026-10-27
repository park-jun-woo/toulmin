//ff:func feature=tangl type=util control=sequence
//ff:what TestReadStatement — readStatement consumes one statement and rejects malformed or body-less ones
package stream

import (
	"testing"

	"github.com/geul-org/gogeul/pkg/toulmin"

	"github.com/geul-org/gogeul/pkg/core"
)

// TestReadStatement reads the Subject statement of an encoded document (TEXT,
// NAMESPACE, LABEL — three records starting after SOURCE, SECTION and its
// heading PROSE), then checks a malformed NODE and a prefix with no body.
func TestReadStatement(t *testing.T) {
	b, err := WriteDoc(buildTransferDoc())
	if err != nil {
		t.Fatalf("WriteDoc: %v", err)
	}
	d, packets := docPackets(t, b)
	next, err := readStatement(d, packets, 3)
	if err != nil || next != 6 || len(d.blocks) != 1 || d.ast.Subject == "" {
		t.Fatalf("success: next=%d blocks=%d subject=%q err=%v", next, len(d.blocks), d.ast.Subject, err)
	}
	badNode := rawPacket{Class: classNode, Rec: core.Record{Kind: core.KindEdge, Codec: toulmin.CodecCaseDeclare}}
	if _, err := readStatement(newDec(testWire), []rawPacket{badNode}, 0); err == nil {
		t.Fatal("malformed NODE: want error")
	}
	if _, err := readStatement(newDec(testWire), []rawPacket{{Class: classPayload}}, 0); err == nil {
		t.Fatal("prefix with no body: want error")
	}
}
