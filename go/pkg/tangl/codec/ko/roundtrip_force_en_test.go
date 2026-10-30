//ff:func feature=tangl type=parser control=sequence
//ff:what TestRoundtripForceEn — EN attack-force fixture: .geul mode-1 byte identity, EN mode-2 identity, KO cross-render re-encode
package ko

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangeul"
	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/codec/en"
)

// TestRoundtripForceEn pins the attack-force grammar for an EN-authored document
// in canonical form (`don't X when Y` beside `Y attacks X with force F`): Save→Load
// keeps mode-1 bytes identical, the loaded AST equals the encoded one, EN mode-2
// reproduces the canonical source byte for byte, and the KO cross-render
// (`decode --locale ko`) re-encodes through the KO codec to the same AST — the
// force survives both locales.
func TestRoundtripForceEn(t *testing.T) {
	src, err := os.ReadFile(filepath.Join("testdata", "force-en.md"))
	if err != nil {
		t.Fatal(err)
	}
	d, err := en.ENCodec{}.Encode(src, "force-en.md")
	if err != nil {
		t.Fatalf("Encode: %v", err)
	}
	out := filepath.Join(t.TempDir(), "force.geul")
	if err := tangeul.Save(d, out); err != nil {
		t.Fatalf("Save: %v", err)
	}
	loaded, err := tangeul.Load(out)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if !bytes.Equal(doc.RenderVerbatim(loaded), src) {
		t.Fatal("mode-1 decode is not byte-identical to the EN source")
	}
	loaded.AST.Path = d.AST.Path // Load stamps the .geul path
	if !tangeul.EqualIgnoringLine(d.AST, loaded.AST) {
		t.Fatal("stream round-trip changed the AST")
	}
	enOut, err := doc.Render(loaded, en.ENCodec{})
	if err != nil {
		t.Fatalf("EN Render: %v", err)
	}
	if !bytes.Equal(enOut, src) {
		t.Fatalf("EN mode-2 of a canonical source is not identical\n--- got ---\n%s", enOut)
	}
	koOut, err := doc.Render(loaded, KoCodec{})
	if err != nil {
		t.Fatalf("KO Render: %v", err)
	}
	re, err := KoCodec{}.Encode(koOut, "force-en.md")
	if err != nil {
		t.Fatalf("KO re-encode: %v\n%s", err, koOut)
	}
	if !tangeul.EqualIgnoringLine(d.AST, re.AST) {
		t.Fatalf("KO cross-render re-encodes to a different AST\n%s", koOut)
	}
}
