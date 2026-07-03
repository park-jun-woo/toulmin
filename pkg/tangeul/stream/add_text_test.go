//ff:func feature=tangl type=util control=sequence
//ff:what TestAddText — addText issues ascending TIDs and appends the Text payload words in issue order
package stream

import (
	"reflect"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/packet"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// TestAddText calls addText twice on the same Writer state and checks that the
// TIDs issued are ascending from 0, textNext advances by one each call, and the
// prefix accumulates both Text payloads in issue order, matching the encoding
// EncodeText itself produces.
func TestAddText(t *testing.T) {
	e := &enc{}
	var prefix []word.Word

	tid0 := addText(e, &prefix, "hello")
	if tid0 != 0 {
		t.Errorf("first TID = %d, want 0", tid0)
	}
	if e.textNext != 1 {
		t.Errorf("textNext after first call = %d, want 1", e.textNext)
	}
	want0 := packet.EncodeText(packet.TextPacket{Bytes: []byte("hello")})
	if !reflect.DeepEqual(prefix, want0) {
		t.Errorf("prefix after first call = %v, want %v", prefix, want0)
	}

	tid1 := addText(e, &prefix, "world")
	if tid1 != 1 {
		t.Errorf("second TID = %d, want 1", tid1)
	}
	if e.textNext != 2 {
		t.Errorf("textNext after second call = %d, want 2", e.textNext)
	}
	want1 := append(append([]word.Word(nil), want0...), packet.EncodeText(packet.TextPacket{Bytes: []byte("world")})...)
	if !reflect.DeepEqual(prefix, want1) {
		t.Errorf("prefix after second call = %v, want %v", prefix, want1)
	}
}
