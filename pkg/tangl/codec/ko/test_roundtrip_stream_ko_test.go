//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what TestRoundtripStreamKo — mode-1 byte identity through a full .tangeul serialization (Save then Load)
package ko

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
)

// TestRoundtripStreamKo pins mode-1 byte identity across a full .tangeul
// serialization for the KO codec (plans/tangl-ko §7 / TANGEUL §9): each KO fixture
// is encoded, the Doc is saved to a .tangeul file, loaded back, and RenderVerbatim
// of the loaded Doc must equal the original source byte-for-byte. americano-ko is
// excluded because it does a `실행한다` on the undeclared alias `기록`, which the
// stream Writer rejects at Save (the same reason the EN suite excludes americano);
// it stays an in-memory-only fixture covered by TestFixtureParity.
func TestRoundtripStreamKo(t *testing.T) {
	for _, name := range []string{"access", "transfer", "rules", "rules_expr"} {
		p := filepath.Join("testdata", name+"-ko.md")
		src, err := os.ReadFile(p)
		if err != nil {
			t.Fatalf("read %s: %v", p, err)
		}
		d, err := KoCodec{}.Encode(src, p)
		if err != nil {
			t.Errorf("%s: Encode: %v", name, err)
			continue
		}
		out := filepath.Join(t.TempDir(), "roundtrip.tangeul")
		if err := tangeul.Save(d, out); err != nil {
			t.Errorf("%s: Save: %v", name, err)
			continue
		}
		loaded, err := tangeul.Load(out)
		if err != nil {
			t.Errorf("%s: Load: %v", name, err)
			continue
		}
		if !bytes.Equal(doc.RenderVerbatim(loaded), src) {
			t.Errorf("%s: RenderVerbatim after stream round-trip is not byte-identical to source", name)
		}
	}
}
