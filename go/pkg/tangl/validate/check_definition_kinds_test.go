//ff:func feature=tangl type=validator control=iteration dimension=1
//ff:what TestCheckDefinitionKinds — invalid fields report once; a valid relative file definition needs Document.Path
package validate

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestCheckDefinitionKinds(t *testing.T) {
	defs := []ast.Definition{
		{Name: "ok const", Kind: ast.ConstDef, Value: "1", Line: 1},
		{Name: "abs file", Kind: ast.FileSpecDef, LoadType: "JSON list", FilePath: "/abs/x.json", Line: 2},
		{Name: "rel file", Kind: ast.FileSpecDef, LoadType: "JSON list", FilePath: "x.json", Line: 3},
		{Name: "bad file", Kind: ast.FileSpecDef, FilePath: "x.json", Line: 4},
		{Name: "bad list", Kind: ast.StringListDef, Value: "v", Line: 5},
	}
	errs := checkDefinitionKinds(&ast.Document{Defs: defs})
	want := []string{
		`:3: definition "rel file" uses a relative file path but Document.Path is empty`,
		`:4: definition "bad file" has invalid fields for FileSpecDef`,
		`:5: definition "bad list" has invalid fields for StringListDef`,
	}
	if len(errs) != len(want) {
		t.Fatalf("errs = %v, want %d", errs, len(want))
	}
	for i, w := range want {
		if errs[i].Error() != w {
			t.Errorf("errs[%d] = %q, want %q", i, errs[i].Error(), w)
		}
	}
	if errs := checkDefinitionKinds(&ast.Document{Path: "/d/p.md", Defs: defs[:3]}); len(errs) != 0 {
		t.Fatalf("with Document.Path the relative file is fine: %v", errs)
	}
}
