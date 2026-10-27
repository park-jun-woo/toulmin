//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what controlsOf — test helper: decode a run of framing records back into ControlPackets
package stream

import (
	"testing"

	"github.com/geul-org/gogeul/pkg/toulmin"

	"github.com/geul-org/gogeul/pkg/core"
)

// controlsOf decodes ws — consecutive framing records at testWire — into their
// ControlPackets, failing the test on any malformed record.
func controlsOf(t *testing.T, ws []core.Word) []toulmin.ControlPacket {
	t.Helper()
	var out []toulmin.ControlPacket
	for i := 0; i < len(ws); {
		r, n, err := core.DecodeRecord(ws[i:])
		if err != nil {
			t.Fatalf("controlsOf: %v", err)
		}
		cp, err := toulmin.DecodeControl(r, testWire)
		if err != nil {
			t.Fatalf("controlsOf: %v", err)
		}
		out = append(out, cp)
		i += n
	}
	return out
}
