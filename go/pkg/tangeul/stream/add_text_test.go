//ff:func feature=tangl type=util control=sequence
//ff:what TestAddText — addText issues consecutive Edge Ids and appends the matching TEXT records
package stream

import (
	"reflect"
	"testing"

	"github.com/geul-org/geul/go/pkg/toulmin"

	"github.com/geul-org/geul/go/pkg/core"
)

// TestAddText checks two calls take Edge Ids 1 and 2 and append exactly the TEXT
// records EncodeText produces for them, and that an id past the width fails.
func TestAddText(t *testing.T) {
	e := testEnc(nil)
	var prefix []core.Word
	id0, err0 := addText(e, &prefix, "hello")
	id1, err1 := addText(e, &prefix, "world")
	if err0 != nil || err1 != nil || id0 != 1 || id1 != 2 {
		t.Fatalf("ids %d %d errs %v %v, want 1 2", id0, id1, err0, err1)
	}
	w0, _ := toulmin.EncodeText(toulmin.TextPacket{TID: 1, Bytes: []byte("hello")}, testWire)
	w1, _ := toulmin.EncodeText(toulmin.TextPacket{TID: 2, Bytes: []byte("world")}, testWire)
	if want := append(w0, w1...); !reflect.DeepEqual(prefix, want) {
		t.Fatalf("prefix = %04X, want %04X", prefix, want)
	}
	e.next = 0x10000
	if _, err := addText(e, &prefix, "x"); err == nil {
		t.Fatal("id past the width: want error")
	}
}
