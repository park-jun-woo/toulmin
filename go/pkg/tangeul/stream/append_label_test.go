//ff:func feature=tangl type=util control=sequence
//ff:what TestAppendLabel — appendLabel puts a TEXT in the prefix and a LABEL (src → TEXT) in the body
package stream

import (
	"reflect"
	"testing"

	"github.com/geul-org/geul/go/pkg/toulmin"

	"github.com/geul-org/geul/go/pkg/core"
)

// TestAppendLabel checks the TEXT takes Edge Id 1, the LABEL Edge Id 2 pointing
// src 7 at TEXT 1, and that a TEXT id past the width is an error.
func TestAppendLabel(t *testing.T) {
	var pre, body []core.Word
	span := toulmin.Span{Off: 2, Len: 3}
	if err := appendLabel(testEnc(nil), &pre, &body, 7, "hasBalance", span); err != nil {
		t.Fatalf("appendLabel: %v", err)
	}
	wantPre, _ := toulmin.EncodeText(toulmin.TextPacket{TID: 1, Bytes: []byte("hasBalance")}, testWire)
	wantBody, _ := toulmin.EncodeEdge(toulmin.EdgePacket{ID: 2, Prop: toulmin.PropLabel, Src: 7, DstTID: 1, Span: span}, testWire)
	if !reflect.DeepEqual(pre, wantPre) || !reflect.DeepEqual(body, wantBody) {
		t.Fatalf("pre %04X body %04X", pre, body)
	}
	e := testEnc(nil)
	e.next = 0x10000
	if err := appendLabel(e, &pre, &body, 7, "x", span); err == nil {
		t.Fatal("id past the width: want error")
	}
}
