//ff:func feature=tangl type=util control=sequence
//ff:what TestEmitSource — encodes the Source Text payload verbatim
package stream

import "testing"

// TestEmitSource is a single straight-line assertion: emitSource is a bare
// wrapper around packet.EncodeText with no branch of its own.
func TestEmitSource(t *testing.T) {
	got := emitSource([]byte("hello"))
	if len(got) == 0 {
		t.Errorf("emitSource: expected non-empty words")
	}
}
