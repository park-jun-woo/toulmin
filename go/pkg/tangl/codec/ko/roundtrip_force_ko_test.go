//ff:func feature=tangl type=parser control=sequence
//ff:what TestRoundtripForceKo — KO attack-force fixture: .geul mode-1 byte identity, F_FORCE preserved, EN golden and KO mode-2 re-encode
package ko

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangeul"
	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/codec/en"
)

// TestRoundtripForceKo pins the attack-force grammar end to end for a KO-authored
// document that mixes the omitted form (무효다, force 1) with explicit forces
// written non-canonically ("0.50", "1.0", "0.25 의"): Save→Load keeps mode-1
// bytes identical (the authored spelling lives in SOURCE), the loaded AST equals
// the encoded one (an explicit 1.0 stays an explicit Force, the omitted form stays
// nil), the EN cross-render equals the force-en.md golden byte for byte, and the
// KO mode-2 rendering is canonical ("0.5의", "1의") and re-encodes to the same AST.
func TestRoundtripForceKo(t *testing.T) {
	src, err := os.ReadFile(filepath.Join("testdata", "force-ko.md"))
	if err != nil {
		t.Fatal(err)
	}
	golden, err := os.ReadFile(filepath.Join("testdata", "force-en.md"))
	if err != nil {
		t.Fatal(err)
	}
	d, err := KoCodec{}.Encode(src, "force-ko.md")
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
		t.Fatal("mode-1 decode is not byte-identical to the KO source")
	}
	loaded.AST.Path = d.AST.Path // Load stamps the .geul path
	if !tangeul.EqualIgnoringLine(d.AST, loaded.AST) {
		t.Fatalf("stream round-trip changed the AST: %+v", loaded.AST.Cases[0].Attacks)
	}
	atk := loaded.AST.Cases[0].Attacks
	if atk[0].Force != nil || atk[1].Force == nil || *atk[1].Force != 0.5 || atk[2].Force == nil || *atk[2].Force != 1 {
		t.Fatalf("forces = %v %v %v, want nil 0.5 1", atk[0].Force, atk[1].Force, atk[2].Force)
	}
	enOut, err := doc.Render(loaded, en.ENCodec{})
	if err != nil {
		t.Fatalf("EN Render: %v", err)
	}
	if !bytes.Equal(enOut, golden) {
		t.Fatalf("EN cross-render differs from golden\n--- got ---\n%s", enOut)
	}
	koOut, err := doc.Render(loaded, KoCodec{})
	if err != nil {
		t.Fatalf("KO Render: %v", err)
	}
	if !strings.Contains(string(koOut), "`신용 충분` 을 0.5의 세기로 공격한다") || !strings.Contains(string(koOut), "`연체 있음` 을 1의 세기로 공격한다") {
		t.Fatalf("KO mode-2 is not canonical:\n%s", koOut)
	}
	re, err := KoCodec{}.Encode(koOut, "force-ko.md")
	if err != nil {
		t.Fatalf("KO re-encode: %v", err)
	}
	if !tangeul.EqualIgnoringLine(d.AST, re.AST) {
		t.Fatal("KO mode-2 rendering re-encodes to a different AST")
	}
}
