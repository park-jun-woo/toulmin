//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what TestEncodeCorpus — encode the whole fixture corpus and assert tiling, statement count, and byte identity
package en

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// TestEncodeCorpus encodes each parser testdata file and each valid test/ example
// and checks the three encoder invariants: the Blocks pass doc.VerifyTiling, the
// number of StmtBlocks equals the count of AST terminal elements (README decision
// 6), and RenderVerbatim reproduces the source byte-for-byte (mode 1).
func TestEncodeCorpus(t *testing.T) {
	astTerminals := func(a *ast.Document) int {
		n := 1 // Subject is always present
		n += len(a.Sees)
		for _, d := range a.Defs {
			n += 1 + len(d.Fields)
		}
		n += len(a.Rules)
		for _, c := range a.Cases {
			n += 1 + len(c.Requires) + len(c.Nodes) + len(c.Attacks) + len(c.Execs)
		}
		for _, ep := range a.Provides {
			n += 1 + len(ep.Requires) + len(ep.Runs) + len(ep.Checks)
		}
		for _, in := range a.Internals {
			n += 1 + len(in.Runs) + len(in.Checks)
		}
		return n
	}
	stmtBlocks := func(d *doc.Doc) int {
		n := 0
		for _, b := range d.Blocks {
			if b.Kind == doc.StmtBlock {
				n++
			}
		}
		return n
	}

	td := func(name string) string {
		return filepath.Join("..", "..", "parser", "testdata", name)
	}
	ex := func(name string) string {
		return filepath.Join("..", "..", "..", "..", "test", name)
	}
	paths := []string{
		td("americano.md"), td("access.md"), td("transfer.md"),
		td("rules_expr.md"),
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
		if d.Locale != "en" {
			t.Errorf("%s: Locale = %q, want en", p, d.Locale)
		}
		if err := doc.VerifyTiling(d); err != nil {
			t.Errorf("%s: VerifyTiling: %v", p, err)
		}
		if got, want := stmtBlocks(d), astTerminals(d.AST); got != want {
			t.Errorf("%s: StmtBlocks = %d, AST terminals = %d", p, got, want)
		}
		if !bytes.Equal(doc.RenderVerbatim(d), src) {
			t.Errorf("%s: RenderVerbatim is not byte-identical to source", p)
		}
	}
}
