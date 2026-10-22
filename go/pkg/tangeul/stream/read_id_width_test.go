//ff:func feature=tangl type=util control=sequence
//ff:what TestReadIDWidth — reads EE; a missing word or reserved bits fail
package stream

import (
	"testing"

	"github.com/geul-org/geul/go/pkg/core"
)

func TestReadIDWidth(t *testing.T) {
	h := &headInfo{}
	if err := readIDWidth(h, core.Record{Payload: []core.Word{0xC000}}); err != nil || h.idWidth != 8 {
		t.Fatalf("got %d %v", h.idWidth, err)
	}
	if err := readIDWidth(h, core.Record{}); err == nil {
		t.Fatal("missing word: want error")
	}
	if err := readIDWidth(h, core.Record{Payload: []core.Word{0x0001}}); err == nil {
		t.Fatal("reserved bits: want error")
	}
}
