//ff:func feature=tangl type=engine control=sequence
//ff:what TestListsLoadErrorAndSkip — a loader failure aborts with nil; const/struct definitions are not materialized
package materialize

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestListsLoadErrorAndSkip(t *testing.T) {
	got, err := Lists(&ast.Document{Defs: []ast.Definition{{Name: "c", Kind: ast.ConstDef, Value: "1"}, {Name: "s", Kind: ast.StructDef}}}, nil)
	if err != nil || len(got) != 0 {
		t.Fatalf("non-list definitions = %#v %v, want empty", got, err)
	}
	doc := &ast.Document{Path: "/d/p.md", Defs: []ast.Definition{{Name: "f", Kind: ast.FileSpecDef, LoadType: "JSON list", FilePath: "f.json"}}}
	got, err = Lists(doc, testRegistry{&testLoader{fail: true}})
	if err == nil || got != nil {
		t.Fatalf("loader failure = %#v %v, want nil map and error", got, err)
	}
}
