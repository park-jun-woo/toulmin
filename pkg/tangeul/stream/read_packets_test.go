//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestReadPackets — readPackets frames a word stream into packets, empty and multi-packet, and errors on a bad prefix or truncation
package stream

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// TestReadPackets drives readPackets directly over hand-framed word slices: no
// packets (loop runs zero times), one packet, several packets (loop runs
// multiple times), a frame whose prefix nibble is not the Free area (surfaced
// from word.DecodeFrame), and a frame whose declared Len runs past the end of
// the stream (the truncation check local to readPackets).
func TestReadPackets(t *testing.T) {
	ctrl := func(body ...word.Word) []word.Word {
		return append(word.EncodeFrame(word.Frame{Class: word.ClassControl, Len: uint32(len(body))}), body...)
	}

	cases := []struct {
		name    string
		in      []word.Word
		wantErr bool
		wantN   int
	}{
		{"empty stream", nil, false, 0},
		{"single packet", ctrl(0x0001), false, 1},
		{"multiple packets", append(ctrl(0x0001), ctrl(0x0002, 0x0003)...), false, 2},
		{"bad frame prefix", []word.Word{0x0C01, 0x0000}, true, 0},
		{"truncated final packet", []word.Word{word.Word(word.EncodeFrame(word.Frame{Class: word.ClassControl, Len: 3})[0]), 0x0000}, true, 0},
	}

	checkClass := func(t *testing.T, p rawPacket) {
		if p.Class != word.ClassControl {
			t.Errorf("packet class = %d, want ClassControl", p.Class)
		}
	}
	check := func(t *testing.T, in []word.Word, wantErr bool, wantN int) {
		out, err := readPackets(in)
		if wantErr && err == nil {
			t.Fatalf("readPackets(%v): expected error, got nil", in)
		}
		if wantErr {
			return
		}
		if err != nil {
			t.Fatalf("readPackets(%v): unexpected error: %v", in, err)
		}
		if len(out) != wantN {
			t.Fatalf("readPackets(%v): got %d packets, want %d", in, len(out), wantN)
		}
		for _, p := range out {
			checkClass(t, p)
		}
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) { check(t, c.in, c.wantErr, c.wantN) })
	}
}
