//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestReadStatement — success appends a Block and advances next; collectStatement/interpretStmt errors propagate
package stream

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// TestReadStatement drives readStatement directly (white-box) over three cases:
// a real statement's packets sliced out of a written stream (success — appends
// one Block and returns the index just past it), a malformed span-bearing
// packet (collectStatement's spanOfRaw error propagates), and a prefix with no
// following Node/Edge body (interpretStmt's decodeBody "no span-bearing packet"
// error propagates). Together these cover both error-return branches plus the
// success path.
func TestReadStatement(t *testing.T) {
	b, err := WriteDoc(buildTransferDoc())
	if err != nil {
		t.Fatalf("WriteDoc: %v", err)
	}
	ws, err := word.ReadWords(b)
	if err != nil {
		t.Fatalf("ReadWords: %v", err)
	}
	packets, err := readPackets(ws)
	if err != nil {
		t.Fatalf("readPackets: %v", err)
	}

	badNode := rawPacket{Class: word.ClassNode, Words: word.EncodeFrame(word.Frame{Class: word.ClassNode, Len: 0})}
	prefixOnly := rawPacket{Class: word.ClassPayload, Words: word.EncodeFrame(word.Frame{Class: word.ClassPayload, Len: 0})}

	tests := []struct {
		name       string
		packets    []rawPacket
		start      int
		wantErr    bool
		wantNext   int
		wantBlocks int
	}{
		{
			name:       "success",
			packets:    packets,
			start:      5,
			wantErr:    false,
			wantNext:   8,
			wantBlocks: 1,
		},
		{
			name:    "collectStatement error (malformed span-bearing packet)",
			packets: []rawPacket{badNode},
			start:   0,
			wantErr: true,
		},
		{
			name:    "interpretStmt error (no span-bearing packet in body)",
			packets: []rawPacket{prefixOnly},
			start:   0,
			wantErr: true,
		},
	}

	check := func(t *testing.T, packets []rawPacket, start int, wantErr bool, wantNext, wantBlocks int) {
		d := newDec()
		next, err := readStatement(d, packets, start)
		if wantErr && err == nil {
			t.Fatalf("expected error, got nil (next=%d)", next)
		}
		if wantErr {
			return
		}
		if err != nil {
			t.Fatalf("readStatement: %v", err)
		}
		if next != wantNext {
			t.Errorf("next = %d, want %d", next, wantNext)
		}
		if len(d.blocks) != wantBlocks {
			t.Errorf("len(d.blocks) = %d, want %d", len(d.blocks), wantBlocks)
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) { check(t, tt.packets, tt.start, tt.wantErr, tt.wantNext, tt.wantBlocks) })
	}
}
