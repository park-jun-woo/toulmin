package gen

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/pkg/tangl/materialize"
)

type genListLoader struct{}

func (genListLoader) LoadDefinition(materialize.DefinitionLoadRequest) (materialize.DefinitionSnapshot, error) {
	return materialize.DefinitionSnapshot{Strings: []string{"loaded"}}, nil
}

type genListRegistry struct{}

func (genListRegistry) LookupDefinitionLoader(string) (materialize.DefinitionLoader, bool) {
	return genListLoader{}, true
}

func TestGenerateStringListDefinition(t *testing.T) {
	doc := &ast.Document{Subject: "list", Defs: []ast.Definition{{Name: "domains", Kind: ast.StringListDef, Strings: []string{"a,b", "c"}}}}
	got, err := Generate(doc)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(got, `toulmin.NewStringListSpec("domains", []string{"a,b", "c"})`) {
		t.Fatalf("generated list missing:\n%s", got)
	}
}

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
