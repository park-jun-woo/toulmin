//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestReadPackets — readPackets classifies BODY records and rejects records outside the profile
package stream

import (
	"testing"

	"github.com/geul-org/gogeul/pkg/toulmin"

	"github.com/geul-org/gogeul/pkg/core"
)

// TestReadPackets classifies an empty BODY, one record of each class, and
// rejects a Value record and an unassigned TOULMIN codec.
func TestReadPackets(t *testing.T) {
	if out, err := readPackets(nil); err != nil || len(out) != 0 {
		t.Fatalf("empty: %v %v", out, err)
	}
	iss := func(f uint16) core.Record {
		return core.Record{Kind: core.KindEdge, Codec: core.CodecIssuer1Word, Issuer: []core.Word{0}, Format: f}
	}
	body := []core.Record{
		iss(toulmin.FormatSource), iss(toulmin.FormatSection), iss(toulmin.FormatText),
		{Kind: core.KindEdge, Codec: toulmin.CodecPredicateDeclare},
		{Kind: core.KindEdge, Codec: toulmin.CodecRuleDeclare},
		{Kind: core.KindEdge, Codec: toulmin.CodecAttacks},
	}
	want := []recClass{classSource, classControl, classPayload, classPayload, classNode, classLink}
	out, err := readPackets(body)
	if err != nil || len(out) != len(want) {
		t.Fatalf("got %d %v", len(out), err)
	}
	for i := range want {
		if out[i].Class != want[i] {
			t.Fatalf("record %d: class %d want %d", i, out[i].Class, want[i])
		}
	}
	for _, bad := range []core.Record{{Kind: core.KindValue}, {Kind: core.KindEdge, Codec: 0x66}} {
		if _, err := readPackets([]core.Record{bad}); err == nil {
			t.Fatalf("%+v: want error", bad)
		}
	}
}
