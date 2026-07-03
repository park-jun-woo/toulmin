//ff:func feature=tangl type=util control=sequence
//ff:what TestDecodeTextErrors — DecodeText rejects a bad frame, wrong class, truncated bodies and an inconsistent byte length
package packet

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// TestDecodeTextErrors exercises DecodeText's self-contained guards: a
// malformed frame header, a Node-class packet, a body shorter than the
// framed length, a body too short for format+byte length, and a byte length
// that does not match the padded data word count.
func TestDecodeTextErrors(t *testing.T) {
	if _, _, err := DecodeText([]word.Word{0x0000}); err == nil {
		t.Fatalf("bad frame prefix: want error")
	}

	nodeWs, err := EncodeNode(NodePacket{TID: 1, Span: word.Span{TID: 1}})
	if err != nil {
		t.Fatalf("EncodeNode: %v", err)
	}
	if _, _, err := DecodeText(nodeWs); err == nil {
		t.Fatalf("wrong class: want error")
	}

	full := EncodeText(TextPacket{Bytes: []byte("hi")})
	if _, _, err := DecodeText(full[:len(full)-1]); err == nil {
		t.Fatalf("truncated body: want error")
	}

	shortBody := append(word.EncodeFrame(word.Frame{Class: word.ClassPayload, Len: 2}), word.Word(TextFormatUTF8), word.Word(0))
	if _, _, err := DecodeText(shortBody); err == nil {
		t.Fatalf("body too short for format+byte length: want error")
	}

	n := word.EncodeNum32(5)
	inconsistent := append(word.EncodeFrame(word.Frame{Class: word.ClassPayload, Len: 4}), word.Word(TextFormatUTF8), n[0], n[1], word.Word(0x4100))
	if _, _, err := DecodeText(inconsistent); err == nil {
		t.Fatalf("byte length inconsistent with data words: want error")
	}
}
