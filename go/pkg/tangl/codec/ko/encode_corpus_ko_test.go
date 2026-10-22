//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what TestEncodeCorpusKo — StmtBlock count equals the AST terminal count for each KO fixture
package ko

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// TestEncodeCorpusKo checks the README decision-6 invariant on the three KO
// fixtures: the number of StmtBlocks equals the count of AST terminal elements
// (subject + sees + defs/fields + rules + case heads/bodies + endpoint heads/bodies
// + internal heads/bodies), so every statement is anchored exactly once.
func TestEncodeCorpusKo(t *testing.T) {
	astTerminals := func(a *ast.Document) int {
		n := 1 // Subject
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
	for _, name := range []string{"americano", "access", "transfer"} {
		p := filepath.Join("testdata", name+"-ko.md")
		src, err := os.ReadFile(p)
		if err != nil {
			t.Fatalf("read %s: %v", p, err)
		}
		d, err := KoCodec{}.Encode(src, p)
		if err != nil {
			t.Fatalf("%s: Encode: %v", name, err)
		}
		if got, want := stmtBlocks(d), astTerminals(d.AST); got != want {
			t.Errorf("%s: StmtBlocks = %d, AST terminals = %d", name, got, want)
		}
	}
}
