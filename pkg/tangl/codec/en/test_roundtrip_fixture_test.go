//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what TestRoundtripFixtures — L1 byte identity on the struck and CRLF fixtures, plus struck AST-absence
package en

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
)

// TestRoundtripFixtures pins mode-1 byte identity (Phase 004 L1) on the two Phase
// 004 corner-case fixtures: struck.md (a `~~...~~` list item, which the parser
// drops and the scanner lowers to prose) and transfer_crlf.md (CRLF line endings,
// whose `\r` bytes must ride inside the raw-byte line spans). For each,
// RenderVerbatim(Encode(src)) must equal src exactly. It also fixes the struck
// contract: the struck node ("legacy pin") is absent from the AST, so the case
// keeps only its one real node.
func TestRoundtripFixtures(t *testing.T) {
	for _, name := range []string{"struck.md", "transfer_crlf.md"} {
		p := filepath.Join("testdata", name)
		src, err := os.ReadFile(p)
		if err != nil {
			t.Fatalf("read %s: %v", p, err)
		}
		d, err := ENCodec{}.Encode(src, p)
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

	src, err := os.ReadFile(filepath.Join("testdata", "struck.md"))
	if err != nil {
		t.Fatalf("read struck.md: %v", err)
	}
	d, err := ENCodec{}.Encode(src, "struck.md")
	if err != nil {
		t.Fatalf("struck.md: Encode: %v", err)
	}
	if len(d.AST.Cases) != 1 {
		t.Fatalf("struck.md: Cases = %d, want 1", len(d.AST.Cases))
	}
	nodes := d.AST.Cases[0].Nodes
	if len(nodes) != 1 {
		t.Errorf("struck.md: Nodes = %d, want 1 (struck node absent from AST)", len(nodes))
	}
	for _, n := range nodes {
		if n.Name == "legacy pin" {
			t.Errorf("struck.md: struck node %q leaked into the AST", n.Name)
		}
	}
}
