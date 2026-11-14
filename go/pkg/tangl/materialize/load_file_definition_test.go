//ff:func feature=tangl type=engine control=sequence
//ff:what TestLoadFileDefinition — registry, loader lookup, path resolution and loader failures each produce a located error
package materialize

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestLoadFileDefinition(t *testing.T) {
	def := ast.Definition{Name: "f", Kind: ast.FileSpecDef, LoadType: "JSON list", FilePath: "data.json", Line: 4}
	doc := &ast.Document{Path: "/docs/policy.md"}

	if _, err := loadFileDefinition(doc, def, nil); err == nil || err.Error() != `/docs/policy.md:4: definition "f": loader registry is required` {
		t.Fatalf("nil registry: %v", err)
	}
	unknown := def
	unknown.LoadType = "YAML list"
	if _, err := loadFileDefinition(doc, unknown, testRegistry{&testLoader{}}); err == nil || !strings.Contains(err.Error(), `loader "YAML list" is not registered`) {
		t.Fatalf("unregistered loader: %v", err)
	}
	if _, err := loadFileDefinition(doc, def, testRegistry{nil}); err == nil || !strings.Contains(err.Error(), "is not registered") {
		t.Fatalf("nil loader: %v", err)
	}
	if _, err := loadFileDefinition(&ast.Document{}, def, testRegistry{&testLoader{}}); err == nil || !strings.Contains(err.Error(), ":4: definition \"f\": relative file path requires Document.Path") {
		t.Fatalf("unresolvable path: %v", err)
	}
	broken := &testLoader{fail: true}
	if _, err := loadFileDefinition(doc, def, testRegistry{broken}); err == nil || err.Error() != `/docs/policy.md:4: definition "f" (loader "JSON list", path "data.json"): broken` {
		t.Fatalf("loader failure: %v", err)
	}

	l := &testLoader{}
	spec, err := loadFileDefinition(doc, def, testRegistry{l})
	if err != nil {
		t.Fatal(err)
	}
	if spec.SpecName() != "f" || len(spec.StringList()) != 2 || l.request.ResolvedPath != "/docs/data.json" || l.request.SourcePath != "data.json" || l.request.TypeName != "JSON list" || l.request.DefinitionName != "f" {
		t.Fatalf("spec=%#v request=%#v", spec, l.request)
	}
}
