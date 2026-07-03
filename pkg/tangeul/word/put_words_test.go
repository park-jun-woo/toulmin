//ff:func feature=tangl type=util control=sequence
//ff:what TestPutWords — words serialize to big-endian bytes, high byte first
package word

import (
	"bytes"
	"testing"
)

func TestPutWords(t *testing.T) {
	if got := PutWords([]Word{0x1234, 0x00FF}); !bytes.Equal(got, []byte{0x12, 0x34, 0x00, 0xFF}) {
		t.Fatalf("got % X want 12 34 00 FF", got)
	}
}
