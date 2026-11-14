//ff:func feature=tangl type=util control=sequence
//ff:what isEdgeRaw — encode an Edge packet as a rawPacket for interpretStmt tests
package stream

import (
	"testing"

	"github.com/geul-org/gogeul/pkg/toulmin"
)

// isEdgeRaw encodes an Edge packet as a rawPacket for interpretStmt tests.
func isEdgeRaw(t *testing.T, p toulmin.EdgePacket) rawPacket {
	t.Helper()
	if p.ID == 0 {
		p.ID, isEdgeNext = isEdgeNext, isEdgeNext+1
	}
	ws, err := toulmin.EncodeEdge(p, testWire)
	if err != nil {
		t.Fatalf("EncodeEdge: %v", err)
	}
	return rawOf(t, ws)
}

// isEdgeNext numbers records built without an Edge Id so each is unique (Edge Id 0 is
// reserved and every id is registered once by the reader).
var isEdgeNext toulmin.TID = 20000
