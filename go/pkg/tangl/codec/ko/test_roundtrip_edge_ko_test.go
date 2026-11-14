//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what TestRoundtripEdgeKo — mode-1 byte identity on the struck/blank-line and CRLF KO fixtures, plus struck AST-absence
package ko

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
)

// TestRoundtripEdgeKo pins mode-1 byte identity on the KO corner-case fixtures:
// struck-ko.md (a `~~...~~` list item the parser drops and the scanner lowers to
// prose, plus a doubled blank line and a trailing blank line) and
// transfer_crlf-ko.md (CRLF line endings, whose `\r` bytes must ride inside the
// raw-byte line spans). For each, VerifyTiling holds and RenderVerbatim(Encode(src))
// equals src exactly. It also fixes the struck contract: the struck node
// ("레거시 핀") is absent from the AST, so the sole case keeps only its one real node.
func TestRoundtripEdgeKo(t *testing.T) {
	for _, name := range []string{"struck-ko.md", "transfer_crlf-ko.md"} {
		p := filepath.Join("testdata", name)
		src, err := os.ReadFile(p)
		if err != nil {
			t.Fatalf("read %s: %v", p, err)
		}
		d, err := KoCodec{}.Encode(src, p)
		if err != nil {
			t.Errorf("%s: Encode: %v", name, err)
			continue
		}
		if err := doc.VerifyTiling(d); err != nil {
			t.Errorf("%s: VerifyTiling: %v", name, err)
		}
		if !bytes.Equal(doc.RenderVerbatim(d), src) {
			t.Errorf("%s: RenderVerbatim is not byte-identical to source", name)
		}
	}

	src, err := os.ReadFile(filepath.Join("testdata", "struck-ko.md"))
	if err != nil {
		t.Fatalf("read struck-ko.md: %v", err)
	}
	d, err := KoCodec{}.Encode(src, "struck-ko.md")
	if err != nil {
		t.Fatalf("struck-ko.md: Encode: %v", err)
	}
	if len(d.AST.Cases) != 1 {
		t.Fatalf("struck-ko.md: Cases = %d, want 1", len(d.AST.Cases))
	}
	nodes := d.AST.Cases[0].Nodes
	if len(nodes) != 1 {
		t.Errorf("struck-ko.md: Nodes = %d, want 1 (struck node absent from AST)", len(nodes))
	}
	for _, n := range nodes {
		if n.Name == "레거시 핀" {
			t.Errorf("struck-ko.md: struck node %q leaked into the AST", n.Name)
		}
	}
}
