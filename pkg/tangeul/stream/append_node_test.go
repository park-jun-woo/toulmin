//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestAppendNode — a valid Node appends its words to body; an unencodable Node errors and leaves body untouched
package stream

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/packet"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// TestAppendNode covers both branches of appendNode: EncodeNode succeeding (the
// words are appended to *body) and EncodeNode failing because F_FUNC is set but
// Func is nil (the error propagates and body is left untouched).
func TestAppendNode(t *testing.T) {
	cases := []struct {
		name    string
		p       packet.NodePacket
		wantErr bool
	}{
		{"valid node", packet.NodePacket{Type: codebook.NodeRule, Span: word.Span{TID: 1}}, false},
		{"F_FUNC set without Func errors", packet.NodePacket{Flags: codebook.NodeFlagFunc}, true},
	}

	run := func(t *testing.T, i int) {
		c := cases[i]
		var body []word.Word
		err := appendNode(&body, c.p)
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
			t.Fatalf("appendNode: %v", err)
		}
		want, err := packet.EncodeNode(c.p)
		if err != nil {
			t.Fatalf("EncodeNode reference: %v", err)
		}
		if len(body) != len(want) {
			t.Errorf("body = %v, want %v", body, want)
		}
	}

	for i := range cases {
		t.Run(cases[i].name, func(t *testing.T) { run(t, i) })
	}
}
