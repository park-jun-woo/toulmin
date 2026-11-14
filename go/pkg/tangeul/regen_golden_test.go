//ff:func feature=tangl type=util control=sequence
//ff:what regenGolden — (re)write one golden .geul under update, then verify its reload AST
package tangeul

import (
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/parser"
)

// regenGolden handles one fixture for TestGoldenRegen: it parses the .md, and when
// update is set rewrites the golden via assembleDoc -> Save. It then Loads the
// golden and checks its AST equals the text-parse AST up to Line (Document.Path is
// aligned first, since Load stamps the .geul path). A missing golden without
// update is a hard failure that points at TANGEUL_UPDATE=1.
func regenGolden(t *testing.T, name string, update bool) {
	t.Helper()
	md := filepath.Join("..", "tangl", "parser", "testdata", name+".md")
	a, err := parser.Parse(md)
	if err != nil {
		t.Fatalf("parse %s: %v", name, err)
	}
	golden := filepath.Join("testdata", name+".geul")
	if update {
		if err := Save(assembleDoc(name, a), golden); err != nil {
			t.Fatalf("regenerate %s: %v", name, err)
		}
	}
	d, err := Load(golden)
	if err != nil {
		t.Fatalf("load %s (run with TANGEUL_UPDATE=1 to create): %v", name, err)
	}
	d.AST.Path = a.Path
	if !EqualIgnoringLine(a, d.AST) {
		t.Fatalf("%s: golden AST drifted from text-parse AST", name)
	}
}
