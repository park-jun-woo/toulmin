//ff:func feature=tangl type=codegen control=sequence
//ff:what TestGenerateFileDefinitionEmbedsSnapshot — a file definition is embedded as a loaded snapshot, never read at runtime
package gen

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestGenerateFileDefinitionEmbedsSnapshot(t *testing.T) {
	doc := &ast.Document{Path: filepath.Join(t.TempDir(), "p.md"), Subject: "list", Defs: []ast.Definition{{Name: "domains", Kind: ast.FileSpecDef, LoadType: "JSON list", FilePath: "x.json"}}}
	got, err := GenerateWithOptions(doc, GenerateOptions{Loaders: genListRegistry{}})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(got, `toulmin.NewStringListSpec("domains", []string{"loaded"})`) || strings.Contains(got, "ReadFile") {
		t.Fatalf("generated snapshot mismatch:\n%s", got)
	}
}
