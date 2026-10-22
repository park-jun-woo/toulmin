//ff:func feature=tangl type=util control=sequence
//ff:what emitHead — the HEAD Meta records of a TANGEUL document plus the BODY marker
package stream

import (
	"github.com/geul-org/geul/go/pkg/core"
	"github.com/geul-org/geul/go/pkg/toulmin"
)

// emitHead encodes the HEAD of a TANGEUL document (GEUL t002 TOULMIN 문서):
// SPEC-VERSION and PROFILE as UTF8-TEXT Values, ID-WIDTH only when the Edge Id
// width is not the 1-word default, DEFAULT-LANGUAGE only when the authoring
// locale is known, then the BODY marker. The Meta records are a few words, so
// encoding them cannot fail; only an undeclarable idWidth is an error.
func emitHead(locale string, idWidth int) ([]core.Word, error) {
	out, _ := metaRecord(core.MetaSpecVersion, core.EncodeBytesValue(core.CodecUTF8TextEven, []byte(toulmin.SpecVersion)))
	prof, _ := metaRecord(core.MetaProfile, core.EncodeBytesValue(core.CodecUTF8TextEven, []byte(toulmin.ProfileString())))
	out = append(out, prof...)
	if idWidth != 1 {
		code, err := core.WidthCode(idWidth)
		if err != nil {
			return nil, err
		}
		iw, _ := core.EncodeRecord(core.Record{Kind: core.KindMeta, Codec: core.MetaIDWidth, Payload: []core.Word{core.Word(uint16(code) << 14)}})
		out = append(out, iw...)
	}
	if locale != "" {
		lang, _ := metaRecord(core.MetaDefaultLanguage, core.EncodeBytesValue(core.CodecLanguageBCP47Even, []byte(locale)))
		out = append(out, lang...)
	}
	body, _ := core.EncodeRecord(core.Record{Kind: core.KindMeta, Codec: core.MetaBody})
	return append(out, body...), nil
}
