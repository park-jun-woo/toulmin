//ff:func feature=tangl type=parser control=sequence
//ff:what TestEncodeStruckKo — a struck case node is dropped from the AST yet the doc still tiles byte-identically
package ko

import (
	"bytes"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
)

// TestEncodeStruckKo pins the v0.1 struck contract: a `~~...~~` list item is absent
// from the AST (the case keeps only its one real node) but its bytes are covered as
// prose, so the tiling holds and RenderVerbatim reproduces the source exactly.
func TestEncodeStruckKo(t *testing.T) {
	src := []byte("## tangl:주제\n- 이 문서는 `x` 이다\n\n## tangl:사안\n- `c` 의 경우에\n  - `인증` 은 일반 규칙이다\n  - ~~`레거시` 은 일반 규칙이다~~\n")
	d, err := KoCodec{}.Encode(src, "t.md")
	if err != nil {
		t.Fatalf("Encode: %v", err)
	}
	if err := doc.VerifyTiling(d); err != nil {
		t.Fatalf("VerifyTiling: %v", err)
	}
	if !bytes.Equal(doc.RenderVerbatim(d), src) {
		t.Error("RenderVerbatim is not byte-identical to source")
	}
	if len(d.AST.Cases) != 1 || len(d.AST.Cases[0].Nodes) != 1 {
		t.Fatalf("Cases/Nodes = %d/%d, want 1/1 (struck node absent)", len(d.AST.Cases), len(d.AST.Cases[0].Nodes))
	}
	if d.AST.Cases[0].Nodes[0].Name != "인증" {
		t.Errorf("surviving node = %q, want 인증", d.AST.Cases[0].Nodes[0].Name)
	}
}
