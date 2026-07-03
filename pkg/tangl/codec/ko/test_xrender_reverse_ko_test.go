//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what TestXrenderReverseKo — KO-authored fixtures render to EN and re-encode to the same AST (reverse cross-locale mode 2)
package ko

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/pkg/tangl/codec/en"
)

// TestXrenderReverseKo pins the reverse cross-locale mode-2 direction (plans/tangl-ko
// §7): a KO-authored document is encoded by the KO codec (d1), rendered as an EN
// audit view with doc.Render + ENCodec, then re-encoded by the EN codec (d2). The EN
// rendering of a KO source must be re-absorbed to the same graph, so d2's AST equals
// d1's up to Line — KO authoring and EN reading share one AST. This complements
// TestRenderStmtXrenderKo (the EN→KO forward direction).
func TestXrenderReverseKo(t *testing.T) {
	for _, name := range []string{"americano", "access", "transfer", "rules", "quest_gate"} {
		p := filepath.Join("testdata", name+"-ko.md")
		src, err := os.ReadFile(p)
		if err != nil {
			t.Fatalf("read %s: %v", p, err)
		}
		d1, err := KoCodec{}.Encode(src, p)
		if err != nil {
			t.Errorf("%s: KO Encode: %v", name, err)
			continue
		}
		rendered, err := doc.Render(d1, en.ENCodec{})
		if err != nil {
			t.Errorf("%s: EN Render: %v", name, err)
			continue
		}
		d2, err := en.ENCodec{}.Encode(rendered, p)
		if err != nil {
			t.Errorf("%s: EN re-encode failed: %v\n---\n%s", name, err, rendered)
			continue
		}
		if !tangeul.EqualIgnoringLine(d1.AST, d2.AST) {
			t.Errorf("%s: EN-rendered AST differs from KO source\n--- rendered ---\n%s", name, rendered)
		}
	}
}
