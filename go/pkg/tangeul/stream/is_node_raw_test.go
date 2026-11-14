//ff:func feature=tangl type=util control=sequence
//ff:what isNodeRaw — encode a Node packet as a rawPacket for interpretStmt tests
package stream

import (
	"testing"

	"github.com/geul-org/gogeul/pkg/toulmin"
)

// isNodeRaw encodes a Node packet as a rawPacket for interpretStmt tests.
func isNodeRaw(t *testing.T, p toulmin.NodePacket) rawPacket {
	t.Helper()
	if p.TID == 0 {
		p.TID, isNodeNext = isNodeNext, isNodeNext+1
	}
	ws, err := toulmin.EncodeNode(p, testWire)
	if err != nil {
		t.Fatalf("EncodeNode: %v", err)
	}
	return rawOf(t, ws)
}

// isNodeNext numbers records built without an Edge Id so each is unique (Edge Id 0 is
// reserved and every id is registered once by the reader).
var isNodeNext toulmin.TID = 30000
