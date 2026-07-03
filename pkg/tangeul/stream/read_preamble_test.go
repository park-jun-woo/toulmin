//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestReadPreamble — readPreamble accepts a valid header/boundary/source and rejects each malformed variant
package stream

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/packet"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// TestReadPreamble drives readPreamble directly over hand-built []rawPacket
// slices, one subtest per branch: too few packets, a bad STREAM_HEADER (via
// checkVersion), a bad SEGMENT_BOUNDARY (via checkSegmentBoundary), a third
// packet that is not the Source Text payload, a Source payload that fails to
// decode, and the success path.
func TestReadPreamble(t *testing.T) {
	headerAndBoundary, err := readPackets(emitHeader(codebook.LocaleUnknown))
	if err != nil {
		t.Fatalf("readPackets(emitHeader()): %v", err)
	}
	if len(headerAndBoundary) != 2 {
		t.Fatalf("emitHeader() produced %d packets, want 2", len(headerAndBoundary))
	}
	header, boundary := headerAndBoundary[0], headerAndBoundary[1]

	sourceWords := packet.EncodeText(packet.TextPacket{Bytes: []byte("hello")})
	source := rawPacket{Class: word.ClassPayload, Words: sourceWords}

	nonControl := rawPacket{Class: word.ClassPayload, Words: word.EncodeFrame(word.Frame{Class: word.ClassPayload, Len: 0})}
	nonPayload := rawPacket{Class: word.ClassNode, Words: word.EncodeFrame(word.Frame{Class: word.ClassNode, Len: 0})}

	badTextWords := append([]word.Word(nil), sourceWords...)
	badTextWords[1] = 0x9999 // corrupt the format code so packet.DecodeText fails
	badText := rawPacket{Class: word.ClassPayload, Words: badTextWords}

	cases := []struct {
		name    string
		packets []rawPacket
		wantErr bool
	}{
		{"too short", []rawPacket{header, boundary}, true},
		{"bad header (checkVersion error)", []rawPacket{nonControl, boundary, source}, true},
		{"bad boundary (checkSegmentBoundary error)", []rawPacket{header, nonControl, source}, true},
		{"third packet not payload", []rawPacket{header, boundary, nonPayload}, true},
		{"source fails to decode", []rawPacket{header, boundary, badText}, true},
		{"valid preamble", []rawPacket{header, boundary, source}, false},
	}

	check := func(t *testing.T, packets []rawPacket, wantErr bool) {
		src, _, err := readPreamble(packets)
		if wantErr && err == nil {
			t.Fatalf("readPreamble: expected error, got nil (src=%+v)", src)
		}
		if wantErr {
			return
		}
		if err != nil {
			t.Fatalf("readPreamble: unexpected error: %v", err)
		}
		if string(src.Bytes) != "hello" {
			t.Errorf("readPreamble: Source = %q, want %q", src.Bytes, "hello")
		}
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) { check(t, c.packets, c.wantErr) })
	}
}
