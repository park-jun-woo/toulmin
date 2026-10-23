//ff:func feature=tangl type=util control=sequence
//ff:what TestAppendLabel — a TEXT on the prefix and a LABEL on the body; a SPAN past the width records an error
package stream

import (
	"reflect"
	"testing"

	"github.com/geul-org/geul/go/pkg/core"
	"github.com/geul-org/geul/go/pkg/toulmin"
)

func TestAppendLabel(t *testing.T) {
	var pre, body []core.Word
	span := toulmin.Span{Off: 2, Len: 3}
	e := testEnc(nil)
	appendLabel(e, &pre, &body, 7, "hasBalance", span)
	wantPre, _ := toulmin.EncodeText(toulmin.TextPacket{TID: 1, Bytes: []byte("hasBalance")}, testWire)
	wantBody, _ := toulmin.EncodeEdge(toulmin.EdgePacket{ID: 2, Prop: toulmin.PropLabel, Src: 7, DstTID: 1, Span: span}, testWire)
	if e.err != nil || !reflect.DeepEqual(pre, wantPre) || !reflect.DeepEqual(body, wantBody) {
		t.Fatalf("pre %04X body %04X err %v", pre, body, e.err)
	}
	bad := testEnc(nil)
	var bp, bb []core.Word
	appendLabel(bad, &bp, &bb, 7, "x", toulmin.Span{Off: 0x10000})
	if bad.err == nil || len(bb) != 0 {
		t.Fatalf("SPAN past the width: err %v body %04X", bad.err, bb)
	}
}
