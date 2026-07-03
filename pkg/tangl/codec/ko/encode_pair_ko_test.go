//ff:func feature=tangl type=parser control=sequence
//ff:what encodePairKo — encode a fixture's KO and EN sources into their two Docs
package ko

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/pkg/tangl/codec/en"
)

// encodePairKo reads the KO and EN sources for fixture name and encodes each with
// its own codec, returning the two Docs. Any read or Encode failure is fatal — the
// shared fixtures must always absorb — so callers get two non-nil Docs to compare.
func encodePairKo(t *testing.T, name string) (*doc.Doc, *doc.Doc) {
	koSrc, err := os.ReadFile(filepath.Join("testdata", name+"-ko.md"))
	if err != nil {
		t.Fatalf("read %s-ko: %v", name, err)
	}
	enSrc, err := os.ReadFile(filepath.Join("testdata", name+"-en.md"))
	if err != nil {
		t.Fatalf("read %s-en: %v", name, err)
	}
	dKo, err := KoCodec{}.Encode(koSrc, name+"-ko.md")
	if err != nil {
		t.Fatalf("%s: KO Encode: %v", name, err)
	}
	dEn, err := en.ENCodec{}.Encode(enSrc, name+"-en.md")
	if err != nil {
		t.Fatalf("%s: EN Encode: %v", name, err)
	}
	return dKo, dEn
}
