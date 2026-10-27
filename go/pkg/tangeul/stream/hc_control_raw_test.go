//ff:func feature=tangl type=util control=sequence
//ff:what hcControlRaw — encode a Control packet as a rawPacket for handleControl tests
package stream

import (
	"testing"

	"github.com/geul-org/gogeul/pkg/toulmin"
)

// hcControlRaw encodes a Control packet as a rawPacket for handleControl tests.
func hcControlRaw(t *testing.T, cp toulmin.ControlPacket) rawPacket {
	t.Helper()
	if cp.ID == 0 {
		cp.ID = 1
	}
	ws, err := toulmin.EncodeControl(cp, testWire)
	if err != nil {
		t.Fatalf("EncodeControl: %v", err)
	}
	return rawOf(t, ws)
}
