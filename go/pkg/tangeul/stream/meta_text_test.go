//ff:func feature=tangl type=util control=sequence
//ff:what TestMetaText — a single text Value payload is read; trailing words, bad records and wrong codecs fail
package stream

import (
	"testing"

	"github.com/geul-org/gogeul/pkg/core"
)

// TestMetaText reads the text of a well-formed Meta and rejects a payload with
// words after its Value, an undecodable payload and a non-text Value.
func TestMetaText(t *testing.T) {
	v, _ := core.EncodeRecord(core.EncodeBytesValue(core.CodecUTF8TextEven, []byte("hi")))
	if s, err := metaText(core.Record{Payload: v}, core.CodecUTF8TextEven); err != nil || s != "hi" {
		t.Fatalf("got %q %v", s, err)
	}
	if _, err := metaText(core.Record{Payload: append(v, 0)}, core.CodecUTF8TextEven); err == nil {
		t.Fatal("trailing word: want error")
	}
	if _, err := metaText(core.Record{Payload: []core.Word{0xE000}}, core.CodecUTF8TextEven); err == nil {
		t.Fatal("bad record: want error")
	}
	if _, err := metaText(core.Record{Payload: v}, core.CodecLanguageBCP47Even); err == nil {
		t.Fatal("wrong codec: want error")
	}
}
