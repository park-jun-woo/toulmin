//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what TestRenderStmtCorpusKo — encode each KO fixture, render mode 2, re-encode, and assert the AST is unchanged
package ko

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangeul"
	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
)

// TestRenderStmtCorpusKo is the re-encoding self-consistency judge (Phase 005 §3):
// for every KO fixture it encodes the source, renders it back with doc.Render (which
// calls RenderStmt for every StmtBlock), re-encodes the rendered document, and
// asserts the resulting AST equals the original up to Line. This exercises every
// statement kind and 조사 allomorph in the corpus through the whole codec at once —
// the normalized rendering must be re-absorbed to the same graph.
func TestRenderStmtCorpusKo(t *testing.T) {
	for _, name := range []string{"americano", "access", "transfer", "rules"} {
		p := filepath.Join("testdata", name+"-ko.md")
		src, err := os.ReadFile(p)
		if err != nil {
			t.Fatalf("read %s: %v", p, err)
		}
		d1, err := KoCodec{}.Encode(src, p)
		if err != nil {
			t.Errorf("%s: Encode: %v", name, err)
			continue
		}
		rendered, err := doc.Render(d1, KoCodec{})
		if err != nil {
			t.Errorf("%s: Render: %v", name, err)
			continue
		}
		d2, err := KoCodec{}.Encode(rendered, p)
		if err != nil {
			t.Errorf("%s: re-encode of rendered output failed: %v\n---\n%s", name, err, rendered)
			continue
		}
		if !tangeul.EqualIgnoringLine(d1.AST, d2.AST) {
			t.Errorf("%s: re-encoded AST differs from original\n--- rendered ---\n%s", name, rendered)
		}
	}
}
