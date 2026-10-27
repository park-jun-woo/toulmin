//ff:func feature=tangl type=util control=sequence
//ff:what rawOf — test helper: decode and classify one encoded record into a rawPacket
package stream

import (
	"testing"

	"github.com/geul-org/gogeul/pkg/core"
)

// rawOf decodes the single record ws and classifies it, as readPackets would.
func rawOf(t *testing.T, ws []core.Word) rawPacket {
	t.Helper()
	r, n, err := core.DecodeRecord(ws)
	if err != nil || n != len(ws) {
		t.Fatalf("rawOf: %v (used %d of %d words)", err, n, len(ws))
	}
	c, err := classifyRecord(r)
	if err != nil {
		t.Fatalf("rawOf: %v", err)
	}
	return rawPacket{Class: c, Rec: r}
}
