//ff:func feature=tangl type=util control=sequence
//ff:what TestAddText — addText issues sequential Edge Ids and appends TEXT records to the prefix
package stream

import (
	"reflect"
	"testing"

	"github.com/geul-org/gogeul/pkg/core"
	"github.com/geul-org/gogeul/pkg/toulmin"
)

func TestAddText(t *testing.T) {
	e := testEnc(nil)
	var prefix []core.Word
	id0 := addText(e, &prefix, "hello")
	id1 := addText(e, &prefix, "world")
	if id0 != 1 || id1 != 2 || e.err != nil {
		t.Fatalf("ids %d %d err %v, want 1 2", id0, id1, e.err)
	}
	w0, _ := toulmin.EncodeText(toulmin.TextPacket{TID: 1, Bytes: []byte("hello")}, testWire)
	w1, _ := toulmin.EncodeText(toulmin.TextPacket{TID: 2, Bytes: []byte("world")}, testWire)
	if want := append(w0, w1...); !reflect.DeepEqual(prefix, want) {
		t.Fatalf("prefix = %04X, want %04X", prefix, want)
	}
}
