//ff:func feature=tangl type=util control=sequence
//ff:what TestExpectMetaText — matching text passes; a different text or a bad payload fails
package stream

import (
	"testing"

	"github.com/geul-org/geul/go/pkg/core"
)

func TestExpectMetaText(t *testing.T) {
	v, _ := core.EncodeRecord(core.EncodeBytesValue(core.CodecUTF8TextEven, []byte("geul-t002")))
	r := core.Record{Kind: core.KindMeta, Codec: core.MetaSpecVersion, Payload: v}
	if err := expectMetaText(r, core.CodecUTF8TextEven, "geul-t002"); err != nil {
		t.Fatal(err)
	}
	if err := expectMetaText(r, core.CodecUTF8TextEven, "geul-t003"); err == nil {
		t.Fatal("different text: want error")
	}
	if err := expectMetaText(core.Record{}, core.CodecUTF8TextEven, "x"); err == nil {
		t.Fatal("empty payload: want error")
	}
}
