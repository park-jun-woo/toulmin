//ff:func feature=tangl type=util control=sequence
//ff:what isNodeRaw — encode a Node packet as a rawPacket for interpretStmt tests
package stream

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/packet"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// isNodeRaw encodes a Node packet as a rawPacket for interpretStmt tests.
func isNodeRaw(t *testing.T, p packet.NodePacket) rawPacket {
	t.Helper()
	ws, err := packet.EncodeNode(p)
	if err != nil {
		t.Fatalf("EncodeNode: %v", err)
	}
	return rawPacket{Class: word.ClassNode, Words: ws}
}
