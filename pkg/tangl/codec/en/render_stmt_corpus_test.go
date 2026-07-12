//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what TestRenderStmtCorpus — encode each fixture, render mode 2, re-parse, and assert the AST is unchanged
package en

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/pkg/tangl/parser"
)

// TestRenderStmtCorpus is the re-parse self-consistency judge (Phase 003 §4): for
// every fixture it encodes the source, renders it back with doc.Render (which
// calls RenderStmt for every StmtBlock), re-parses the rendered document, and
// asserts the resulting AST equals the original up to Line. This exercises all
// ~25 statement kinds present in the corpus through the whole codec at once.
func TestRenderStmtCorpus(t *testing.T) {
	td := func(name string) string {
		return filepath.Join("..", "..", "parser", "testdata", name)
	}
	ex := func(name string) string {
		return filepath.Join("..", "..", "..", "..", "test", name)
	}
	paths := []string{
		td("americano.md"), td("access.md"), td("transfer.md"),
		td("supersede_level.md"), td("rules_expr.md"),
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
		rendered, err := doc.Render(d, ENCodec{})
		if err != nil {
			t.Errorf("%s: Render: %v", p, err)
			continue
		}
		reparsed, err := parser.ParseSource(string(rendered), d.AST.Path)
		if err != nil {
			t.Errorf("%s: re-parse of rendered output failed: %v\n---\n%s", p, err, rendered)
			continue
		}
		if !tangeul.EqualIgnoringLine(d.AST, reparsed) {
			t.Errorf("%s: re-parsed AST differs from original\n--- rendered ---\n%s", p, rendered)
		}
	}
}
