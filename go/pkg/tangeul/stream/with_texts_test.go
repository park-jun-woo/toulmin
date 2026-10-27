//ff:func feature=tangl type=util control=sequence
//ff:what TestWithTexts — withTexts issues one TEXT per term in order; no terms issue nothing
package stream

import (
	"reflect"
	"testing"

	"github.com/geul-org/gogeul/pkg/core"
	"github.com/geul-org/gogeul/pkg/toulmin"
)

func TestWithTexts(t *testing.T) {
	var pre []core.Word
	if got := withTexts(testEnc(nil), &pre, nil); len(got) != 0 || len(pre) != 0 {
		t.Fatalf("no terms: %v %v", got, pre)
	}
	if got := withTexts(testEnc(nil), &pre, []string{"amount", "account"}); !reflect.DeepEqual(got, []toulmin.TID{1, 2}) {
		t.Fatalf("ids = %v, want [1 2]", got)
	}
	w0, _ := toulmin.EncodeText(toulmin.TextPacket{TID: 1, Bytes: []byte("amount")}, testWire)
	w1, _ := toulmin.EncodeText(toulmin.TextPacket{TID: 2, Bytes: []byte("account")}, testWire)
	if !reflect.DeepEqual(pre, append(w0, w1...)) {
		t.Fatalf("prefix = %04X", pre)
	}
}
