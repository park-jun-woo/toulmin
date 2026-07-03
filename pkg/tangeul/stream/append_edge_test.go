//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestAppendEdge — a valid Edge appends its words to body; an unencodable Edge errors and leaves body untouched
package stream

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/packet"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// TestAppendEdge covers both branches of appendEdge: EncodeEdge succeeding (the
// words are appended to *body) and EncodeEdge failing on an unknown PropCode
// (the error propagates and body is left as it was).
func TestAppendEdge(t *testing.T) {
	cases := []struct {
		name    string
		ep      packet.EdgePacket
		wantErr bool
	}{
		{"valid edge", packet.EdgePacket{Prop: codebook.PropAttacks, Src: 1, DstTID: 2, Span: word.Span{TID: 1}}, false},
		{"unknown prop code errors", packet.EdgePacket{Prop: codebook.PropCode(0x9999)}, true},
	}

	run := func(t *testing.T, i int) {
		c := cases[i]
		var body []word.Word
		err := appendEdge(&body, c.ep)
		if c.wantErr {
			if err == nil {
				t.Fatalf("want error, got nil")
			}
			if len(body) != 0 {
				t.Errorf("body should stay empty on error, got %v", body)
			}
			return
		}
		if err != nil {
			t.Fatalf("appendEdge: %v", err)
		}
		want, err := packet.EncodeEdge(c.ep)
		if err != nil {
			t.Fatalf("EncodeEdge reference: %v", err)
		}
		if len(body) != len(want) {
			t.Errorf("body = %v, want %v", body, want)
		}
	}

	for i := range cases {
		t.Run(cases[i].name, func(t *testing.T) { run(t, i) })
	}
}
