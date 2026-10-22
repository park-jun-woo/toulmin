//ff:func feature=tangl type=util control=selection
//ff:what readHeadRecord — apply one HEAD Meta record (SPEC-VERSION, PROFILE, ID-WIDTH, DEFAULT-LANGUAGE)
package stream

import (
	"fmt"

	"github.com/geul-org/geul/go/pkg/toulmin"

	"github.com/geul-org/geul/go/pkg/core"
)

// readHeadRecord applies one HEAD Meta record to h. Declaration-class Meta that
// this reader does not know must not be ignored (it may change interpretation),
// so it is an error; source, lineage, verification and policy Meta (codecs
// 0x20-0x5F) only inform trust and handling and are skipped. IMPORT is not part
// of the TANGEUL profile yet.
func readHeadRecord(h *headInfo, r core.Record) error {
	switch r.Codec {
	case core.MetaSpecVersion:
		h.spec = true
		return expectMetaText(r, core.CodecUTF8TextEven, toulmin.SpecVersion)
	case core.MetaProfile:
		h.profile = true
		s, err := metaText(r, core.CodecUTF8TextEven)
		if err != nil {
			return err
		}
		return toulmin.CheckProfile(s)
	case core.MetaIDWidth:
		return readIDWidth(h, r)
	case core.MetaDefaultLanguage:
		s, err := metaText(r, core.CodecLanguageBCP47Even)
		h.locale = s
		return err
	default:
		if r.Codec >= 0x20 && r.Codec < 0x60 {
			return nil
		}
		return fmt.Errorf("tangeul/stream: unsupported HEAD Meta 0x%02X", r.Codec)
	}
}
