//ff:func feature=tangl type=engine control=iteration dimension=1
//ff:what TestMaterializeDefinition — lists and file definitions yield Specs; other kinds are skipped; load errors propagate
package materialize

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestMaterializeDefinition(t *testing.T) {
	doc := &ast.Document{Path: "/d/p.md"}
	spec, ok, err := materializeDefinition(doc, ast.Definition{Name: "l", Kind: ast.StringListDef, Strings: []string{"x"}}, nil)
	if err != nil || !ok || spec.SpecName() != "l" || len(spec.StringList()) != 1 {
		t.Fatalf("string list = %#v %v %v", spec, ok, err)
	}
	spec, ok, err = materializeDefinition(doc, ast.Definition{Name: "f", Kind: ast.FileSpecDef, LoadType: "JSON list", FilePath: "f.json"}, testRegistry{&testLoader{}})
	if err != nil || !ok || spec.SpecName() != "f" || len(spec.StringList()) != 2 {
		t.Fatalf("file = %#v %v %v", spec, ok, err)
	}
	if _, ok, err := materializeDefinition(doc, ast.Definition{Name: "f", Kind: ast.FileSpecDef}, nil); ok || err == nil {
		t.Fatalf("file without registry = ok %v err %v", ok, err)
	}
	for _, kind := range []ast.DefKind{ast.ConstDef, ast.StructDef} {
		if _, ok, err := materializeDefinition(doc, ast.Definition{Name: "c", Kind: kind}, nil); ok || err != nil {
			t.Fatalf("%v must be skipped: ok %v err %v", kind, ok, err)
		}
	}
}
