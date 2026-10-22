//ff:func feature=tangl type=util control=sequence
//ff:what TestWithTexts — withTexts issues one TEXT per term in order and propagates an encode error
package stream

import (
	"reflect"
	"testing"

	"github.com/geul-org/geul/go/pkg/toulmin"

	"github.com/geul-org/geul/go/pkg/core"
)

// TestWithTexts checks no terms yields no records, two terms yield ids 1 and 2
// with their TEXT records, and an id past the width is an error.
func TestWithTexts(t *testing.T) {
	var pre []core.Word
	if got, err := withTexts(testEnc(nil), &pre, nil); err != nil || len(got) != 0 || len(pre) != 0 {
		t.Fatalf("no terms: %v %v %v", got, pre, err)
	}
	got, err := withTexts(testEnc(nil), &pre, []string{"amount", "account"})
	if err != nil || !reflect.DeepEqual(got, []toulmin.TID{1, 2}) {
		t.Fatalf("ids = %v %v, want [1 2]", got, err)
	}
	w0, _ := toulmin.EncodeText(toulmin.TextPacket{TID: 1, Bytes: []byte("amount")}, testWire)
	w1, _ := toulmin.EncodeText(toulmin.TextPacket{TID: 2, Bytes: []byte("account")}, testWire)
	if !reflect.DeepEqual(pre, append(w0, w1...)) {
		t.Fatalf("prefix = %04X", pre)
	}
	e := testEnc(nil)
	e.next = 0x10000
	if _, err := withTexts(e, &pre, []string{"x"}); err == nil {
		t.Fatal("id past the width: want error")
	}
}
