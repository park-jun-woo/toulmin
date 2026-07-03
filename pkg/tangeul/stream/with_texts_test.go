//ff:func feature=tangl type=util control=sequence
//ff:what TestWithTexts — withTexts issues one Text payload and TID per `with` term, in order, and none for an empty term list
package stream

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/packet"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// TestWithTexts exercises the for-range over terms directly: no terms (zero
// iterations, empty result and untouched prefix) and multiple terms (several
// iterations, ascending TIDs and one Text payload per term appended in order).
func TestWithTexts(t *testing.T) {
	t.Run("no terms", func(t *testing.T) {
		e := &enc{}
		var pre []word.Word
		got := withTexts(e, &pre, nil)
		if len(got) != 0 {
			t.Errorf("TIDs = %v, want empty", got)
		}
		if len(pre) != 0 {
			t.Errorf("prefix = %v, want empty", pre)
		}
		if e.textNext != 0 {
			t.Errorf("textNext = %d, want 0", e.textNext)
		}
	})

	t.Run("multiple terms", func(t *testing.T) {
		e := &enc{}
		var pre []word.Word
		got := withTexts(e, &pre, []string{"amount", "account"})

		wantTIDs := []packet.TID{0, 1}
		if len(got) != len(wantTIDs) || got[0] != wantTIDs[0] || got[1] != wantTIDs[1] {
			t.Errorf("TIDs = %v, want %v", got, wantTIDs)
		}
		if e.textNext != 2 {
			t.Errorf("textNext = %d, want 2", e.textNext)
		}

		want := append(
			append([]word.Word(nil), packet.EncodeText(packet.TextPacket{Bytes: []byte("amount")})...),
			packet.EncodeText(packet.TextPacket{Bytes: []byte("account")})...,
		)
		if len(pre) != len(want) {
			t.Fatalf("prefix len = %d, want %d", len(pre), len(want))
		}
		for i := range want {
			if pre[i] != want[i] {
				t.Errorf("prefix[%d] = %v, want %v", i, pre[i], want[i])
			}
		}
	})
}
