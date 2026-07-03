//ff:func feature=tangl type=util control=sequence
//ff:what hcControlRaw — encode a Control packet as a rawPacket for handleControl tests
package stream

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/packet"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// hcControlRaw encodes a Control packet as a rawPacket for handleControl tests.
func hcControlRaw(t *testing.T, cp packet.ControlPacket) rawPacket {
	t.Helper()
	ws, err := packet.EncodeControl(cp)
	if err != nil {
		t.Fatalf("EncodeControl: %v", err)
	}
	return rawPacket{Class: word.ClassControl, Words: ws}
}
