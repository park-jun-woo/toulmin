//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestAppendNode — a valid Node appends its words to body; an unencodable Node errors and leaves body untouched
package stream

import (
	"testing"

	"github.com/geul-org/geul/go/pkg/toulmin"

	"github.com/geul-org/geul/go/pkg/core"
)

// TestAppendNode covers both branches of appendNode: EncodeNode succeeding (the
// words are appended to *body) and EncodeNode failing because F_FUNC is set but
// Func is nil (the error propagates and body is left untouched).
func TestAppendNode(t *testing.T) {
	cases := []struct {
		name    string
		p       toulmin.NodePacket
		wantErr bool
	}{
		{"valid node", toulmin.NodePacket{Type: toulmin.NodeRule, Span: toulmin.Span{}}, false},
		{"F_FUNC set without Func errors", toulmin.NodePacket{Flags: toulmin.NodeFlagFunc}, true},
	}

	run := func(t *testing.T, i int) {
		c := cases[i]
		var body []core.Word
		err := appendNode(testEnc(nil), &body, c.p)
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
		want, err := toulmin.EncodeNode(c.p, testWire)
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
