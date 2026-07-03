//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what TestRoundtripStream — L2 byte identity through the .tangeul stream (Save then Load)
package en

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
)

// TestRoundtripStream pins mode-1 byte identity across a full .tangeul serialization
// (Phase 004 L2): for each fixture it encodes the source, saves the Doc to a
// .tangeul file, loads it back, and asserts RenderVerbatim of the loaded Doc equals
// the original source. americano is excluded because it does a `do` on the
// undeclared alias `log`, which the stream Writer rejects at Save (pinned by
// tangeul.TestAmericanoExcluded); it stays an L1-only fixture.
func TestRoundtripStream(t *testing.T) {
	td := func(name string) string {
		return filepath.Join("..", "..", "parser", "testdata", name)
	}
	ex := func(name string) string {
		return filepath.Join("..", "..", "..", "..", "test", name)
	}
	paths := []string{
		td("access.md"), td("transfer.md"),
		ex("library-renewal.md"), ex("newsletter-tick-warn.md"),
		ex("newsletter-tick-fixed.md"), ex("order-saga.md"),
		ex("factory-line-cascade.md"),
	}
	for _, p := range paths {
		src, err := os.ReadFile(p)
		if err != nil {
			t.Fatalf("read %s: %v", p, err)
		}
		d, err := ENCodec{}.Encode(src, p)
		if err != nil {
			t.Errorf("%s: Encode: %v", p, err)
			continue
		}
		out := filepath.Join(t.TempDir(), "roundtrip.tangeul")
		if err := tangeul.Save(d, out); err != nil {
			t.Errorf("%s: Save: %v", p, err)
			continue
		}
		loaded, err := tangeul.Load(out)
		if err != nil {
			t.Errorf("%s: Load: %v", p, err)
			continue
		}
		if !bytes.Equal(doc.RenderVerbatim(loaded), src) {
			t.Errorf("%s: RenderVerbatim after stream round-trip is not byte-identical to source", p)
		}
	}
}
