//ff:func feature=tangl type=util control=sequence
//ff:what isEdgeRaw — encode an Edge packet as a rawPacket for interpretStmt tests
package stream

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/packet"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// isEdgeRaw encodes an Edge packet as a rawPacket for interpretStmt tests.
func isEdgeRaw(t *testing.T, p packet.EdgePacket) rawPacket {
	t.Helper()
	ws, err := packet.EncodeEdge(p)
	if err != nil {
		t.Fatalf("EncodeEdge: %v", err)
	}
	return rawPacket{Class: word.ClassEdge, Words: ws}
}
