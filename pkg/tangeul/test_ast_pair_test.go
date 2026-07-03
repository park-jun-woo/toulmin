//ff:func feature=tangl type=util control=sequence
//ff:what astPair — load a fixture's AST twice: direct text-parse vs golden .tangeul reload
package tangeul

import (
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/pkg/tangl/parser"
)

// astPair returns the two ASTs the Phase 008 step-1 parity tests compare: the
// "old" path (parser.Parse over the authored .md, the pre-TANGEUL text route) and
// the "new" path (tangeul.Load over the frozen golden .tangeul). Both are returned
// unnormalized — Document.Path differs (.md vs .tangeul) and Line differs (Load
// leaves 0), so callers normalize per their comparison (EqualIgnoringLine for the
// AST, normalizeLoc for validate messages). Parity fixtures are transfer and
// access; americano is excluded (see TestGoldenRegen for the reason).
func astPair(t *testing.T, name string) (old, loaded *ast.Document) {
	t.Helper()
	md := filepath.Join("..", "tangl", "parser", "testdata", name+".md")
	old, err := parser.Parse(md)
	if err != nil {
		t.Fatalf("parse %s: %v", name, err)
	}
	d, err := Load(filepath.Join("testdata", name+".tangeul"))
	if err != nil {
		t.Fatalf("load golden %s: %v", name, err)
	}
	return old, d.AST
}
