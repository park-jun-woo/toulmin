//ff:func feature=tangl type=util control=sequence
//ff:what TestReadPreamble — readPreamble requires SOURCE first and returns its text and SPAN width
package stream

import (
	"testing"

	"github.com/geul-org/gogeul/pkg/toulmin"

	"github.com/geul-org/gogeul/pkg/core"
)

// TestReadPreamble covers an empty BODY, a BODY not opened by SOURCE, a SOURCE
// that fails to decode, and the success path with a 2-word SPAN width.
func TestReadPreamble(t *testing.T) {
	sw, _ := toulmin.EncodeSource(toulmin.TextPacket{TID: 1, Bytes: []byte("hello")}, toulmin.Wire{Span: 2})
	src := rawOf(t, sw)
	if _, _, err := readPreamble(nil); err == nil {
		t.Fatal("empty BODY: want error")
	}
	tw, _ := toulmin.EncodeText(toulmin.TextPacket{TID: 1}, testWire)
	if _, _, err := readPreamble([]rawPacket{rawOf(t, tw)}); err == nil {
		t.Fatal("BODY opened by TEXT: want error")
	}
	bad := rawPacket{Class: classSource, Rec: core.Record{Kind: core.KindEdge, Codec: core.CodecIssuer1Word, Issuer: []core.Word{0}, Format: toulmin.FormatSource}}
	if _, _, err := readPreamble([]rawPacket{bad}); err == nil {
		t.Fatal("malformed SOURCE: want error")
	}
	got, w, err := readPreamble([]rawPacket{src})
	if err != nil || string(got.Bytes) != "hello" || w != 2 {
		t.Fatalf("got %q %d %v", got.Bytes, w, err)
	}
}
