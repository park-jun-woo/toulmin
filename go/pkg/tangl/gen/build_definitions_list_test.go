//ff:func feature=tangl type=codegen control=sequence
//ff:what TestBuildDefinitionsLists — string-list and file definitions render from the snapshot map
package gen

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestBuildDefinitionsLists(t *testing.T) {
	var w strings.Builder
	index := buildDefinitions(&w, []ast.Definition{
		{Name: "inline", Kind: ast.StringListDef, Strings: []string{"ignored"}},
		{Name: "file", Kind: ast.FileSpecDef, LoadType: "JSON list", FilePath: "x.json"},
		{Name: "limit", Kind: ast.ConstDef, Value: "3"},
	}, map[string][]string{"inline": {"a"}, "file": {"b", "c"}})
	out := w.String()
	if !strings.Contains(out, `toulmin.NewStringListSpec("inline", []string{"a"})`) || !strings.Contains(out, `toulmin.NewStringListSpec("file", []string{"b", "c"})`) {
		t.Fatalf("list definitions not rendered from snapshots:\n%s", out)
	}
	if index["inline"].Kind != ast.StringListDef || index["file"].Kind != ast.FileSpecDef || index["file"].Spec != "fileSpec" {
		t.Fatalf("index = %+v", index)
	}
	if _, ok := index["limit"]; !ok {
		t.Fatal("const definition must be indexed too")
	}
}
