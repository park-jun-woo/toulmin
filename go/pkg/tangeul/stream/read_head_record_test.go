//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestReadHeadRecord — HEAD Meta handling: stamps, language, skipped and rejected classes
package stream

import (
	"testing"

	"github.com/geul-org/geul/go/pkg/core"
)

// TestReadHeadRecord applies each HEAD Meta kind: a wrong SPEC-VERSION, an
// unknown PROFILE, a malformed PROFILE payload, a language tag, a skipped
// provenance Meta (AUTHOR) and a rejected removed ID-WIDTH (0x13), unknown
// declaration and IMPORT.
func TestReadHeadRecord(t *testing.T) {
	text := func(codec uint8, even uint8, s string) core.Record {
		v, _ := core.EncodeRecord(core.EncodeBytesValue(even, []byte(s)))
		return core.Record{Kind: core.KindMeta, Codec: codec, Payload: v}
	}
	h := &headInfo{}
	ok := []core.Record{
		text(core.MetaDefaultLanguage, core.CodecLanguageBCP47Even, "ko"),
		{Kind: core.KindMeta, Codec: 0x20},
	}
	for _, r := range ok {
		if err := readHeadRecord(h, r); err != nil {
			t.Fatalf("0x%02X: %v", r.Codec, err)
		}
	}
	if h.locale != "ko" {
		t.Fatalf("head = %+v", h)
	}
	bad := []core.Record{
		core.EncodeSpecVersion(1),
		text(core.MetaProfile, core.CodecUTF8TextEven, "tangeul-0.3"),
		{Kind: core.KindMeta, Codec: core.MetaProfile, Payload: []core.Word{0x0882, 1}},
		{Kind: core.KindMeta, Codec: 0x13, Payload: []core.Word{0x4000}},
		{Kind: core.KindMeta, Codec: 0x12},
		{Kind: core.KindMeta, Codec: 0x60},
	}
	for _, r := range bad {
		if err := readHeadRecord(&headInfo{}, r); err == nil {
			t.Fatalf("0x%02X: want error", r.Codec)
		}
	}
}
