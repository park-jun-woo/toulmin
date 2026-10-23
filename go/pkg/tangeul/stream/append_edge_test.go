//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestAppendEdge — a valid Edge appends its words to body; an unencodable Edge records an error and leaves body untouched
package stream

import (
	"testing"

	"github.com/geul-org/geul/go/pkg/toulmin"

	"github.com/geul-org/geul/go/pkg/core"
)

// TestAppendEdge covers both branches of appendEdge: EncodeEdge succeeding (the
// words are appended to *body) and EncodeEdge failing on an unknown PropCode
// (the error propagates and body is left as it was).
func TestAppendEdge(t *testing.T) {
	cases := []struct {
		name    string
		ep      toulmin.EdgePacket
		wantErr bool
	}{
		{"valid edge", toulmin.EdgePacket{Prop: toulmin.PropAttacks, Src: 1, DstTID: 2, Span: toulmin.Span{}}, false},
		{"unknown prop code errors", toulmin.EdgePacket{Prop: toulmin.PropCode(0x9999)}, true},
	}

	run := func(t *testing.T, i int) {
		c := cases[i]
		var body []core.Word
		e := testEnc(nil)
		appendEdge(e, &body, c.ep)
		if c.wantErr {
			if e.err == nil {
				t.Fatalf("want error, got nil")
			}
			if len(body) != 0 {
				t.Errorf("body should stay empty on error, got %v", body)
			}
			return
		}
		if e.err != nil {
			t.Fatalf("appendEdge: %v", e.err)
		}
		want, err := toulmin.EncodeEdge(c.ep, testWire)
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
