//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestReadDoc — ReadDoc's post-preamble loop: zero statement packets, a mid-stream stepRead failure, and the round-trip already covers success
package stream

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// TestReadDoc exercises the two branches round_trip_test.go, round_trip_rich_test.go
// and read_errors_test.go leave uncovered in ReadDoc's `for i := 3; i <
// len(packets); ` loop: the zero-iteration case (a stream that is only the
// header, boundary and Source — no statement or control packets follow) and the
// stepRead-returns-an-error case (a mid-stream packet that fails, here via an
// unrecognized control code that handleControl rejects).
func TestReadDoc(t *testing.T) {
	full, err := WriteDoc(buildTransferDoc())
	if err != nil {
		t.Fatalf("WriteDoc: %v", err)
	}
	ws, err := word.ReadWords(full)
	if err != nil {
		t.Fatalf("word.ReadWords: %v", err)
	}
	packets, err := readPackets(ws)
	if err != nil {
		t.Fatalf("readPackets: %v", err)
	}
	if len(packets) < 4 {
		t.Fatalf("fixture stream has only %d packets, need a 4th to corrupt", len(packets))
	}
	var preamble []word.Word
	for _, p := range packets[:3] {
		preamble = append(preamble, p.Words...)
	}

	t.Run("zero statement packets after preamble", func(t *testing.T) {
		got, err := ReadDoc(word.PutWords(preamble))
		if err != nil {
			t.Fatalf("ReadDoc: unexpected error: %v", err)
		}
		if got.AST == nil {
			t.Fatalf("ReadDoc: AST is nil")
		}
		if len(got.Blocks) != 0 {
			t.Errorf("ReadDoc: Blocks = %v, want none", got.Blocks)
		}
	})

	t.Run("stepRead error mid-stream", func(t *testing.T) {
		badCtrl := append(word.EncodeFrame(word.Frame{Class: word.ClassControl, Len: 1}), word.Word(0xFF))
		bad := append(append([]word.Word(nil), preamble...), badCtrl...)
		if _, err := ReadDoc(word.PutWords(bad)); err == nil {
			t.Fatalf("ReadDoc: expected stepRead error, got nil")
		}
	})
}
