//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what TestRenderStmtXrenderKo — EN-authored fixtures render to KO and re-encode to the same AST (cross-locale mode 2)
package ko

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/pkg/tangl/codec/en"
)

// TestRenderStmtXrenderKo pins cross-locale mode-2 decoding (README §9 / plans/tangl-ko
// §7): an EN-authored document is encoded by the EN codec (d1), rendered as a KO
// audit view with doc.Render + KoCodec, then re-encoded by the KO codec (d2). The KO
// rendering must be re-absorbed to the same graph as the EN source, so d2's AST
// equals d1's up to Line — EN authoring and KO reading share one AST.
func TestRenderStmtXrenderKo(t *testing.T) {
	for _, name := range []string{"americano", "access", "transfer", "rules"} {
		p := filepath.Join("testdata", name+"-en.md")
		src, err := os.ReadFile(p)
		if err != nil {
			t.Fatalf("read %s: %v", p, err)
		}
		d1, err := en.ENCodec{}.Encode(src, p)
		if err != nil {
			t.Errorf("%s: EN Encode: %v", name, err)
			continue
		}
		rendered, err := doc.Render(d1, KoCodec{})
		if err != nil {
			t.Errorf("%s: KO Render: %v", name, err)
			continue
		}
		d2, err := KoCodec{}.Encode(rendered, p)
		if err != nil {
			t.Errorf("%s: KO re-encode failed: %v\n---\n%s", name, err, rendered)
			continue
		}
		if !tangeul.EqualIgnoringLine(d1.AST, d2.AST) {
			t.Errorf("%s: KO-rendered AST differs from EN source\n--- rendered ---\n%s", name, rendered)
		}
	}
}
