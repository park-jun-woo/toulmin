//ff:func feature=tangl type=validator control=iteration dimension=1
//ff:what TestDefinitionFieldsInvalid — each Kind rejects foreign fields; file definitions need loader and path; unknown kinds fail
package validate

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestDefinitionFieldsInvalid(t *testing.T) {
	ref := &ast.Ref{Name: "NewX"}
	fields := []ast.Field{{Name: "f"}}
	cases := []struct {
		name string
		def  ast.Definition
		want bool
	}{
		{"const value+spec", ast.Definition{Kind: ast.ConstDef, Value: "1", SpecRef: ref}, false},
		{"const with fields", ast.Definition{Kind: ast.ConstDef, Fields: fields}, true},
		{"const with strings", ast.Definition{Kind: ast.ConstDef, Strings: []string{"a"}}, true},
		{"const with loader", ast.Definition{Kind: ast.ConstDef, LoadType: "J"}, true},
		{"const with path", ast.Definition{Kind: ast.ConstDef, FilePath: "x"}, true},
		{"struct fields", ast.Definition{Kind: ast.StructDef, Fields: fields}, false},
		{"struct with value", ast.Definition{Kind: ast.StructDef, Value: "1"}, true},
		{"struct with spec", ast.Definition{Kind: ast.StructDef, SpecRef: ref}, true},
		{"struct with strings", ast.Definition{Kind: ast.StructDef, Strings: []string{"a"}}, true},
		{"struct with loader", ast.Definition{Kind: ast.StructDef, LoadType: "J"}, true},
		{"struct with path", ast.Definition{Kind: ast.StructDef, FilePath: "x"}, true},
		{"list strings", ast.Definition{Kind: ast.StringListDef, Strings: []string{"a"}}, false},
		{"list with value", ast.Definition{Kind: ast.StringListDef, Value: "1"}, true},
		{"list with spec", ast.Definition{Kind: ast.StringListDef, SpecRef: ref}, true},
		{"list with fields", ast.Definition{Kind: ast.StringListDef, Fields: fields}, true},
		{"list with loader", ast.Definition{Kind: ast.StringListDef, LoadType: "J"}, true},
		{"list with path", ast.Definition{Kind: ast.StringListDef, FilePath: "x"}, true},
		{"file loader+path", ast.Definition{Kind: ast.FileSpecDef, LoadType: "J", FilePath: "x"}, false},
		{"file without loader", ast.Definition{Kind: ast.FileSpecDef, FilePath: "x"}, true},
		{"file without path", ast.Definition{Kind: ast.FileSpecDef, LoadType: "J"}, true},
		{"file with value", ast.Definition{Kind: ast.FileSpecDef, LoadType: "J", FilePath: "x", Value: "1"}, true},
		{"file with spec", ast.Definition{Kind: ast.FileSpecDef, LoadType: "J", FilePath: "x", SpecRef: ref}, true},
		{"file with fields", ast.Definition{Kind: ast.FileSpecDef, LoadType: "J", FilePath: "x", Fields: fields}, true},
		{"file with strings", ast.Definition{Kind: ast.FileSpecDef, LoadType: "J", FilePath: "x", Strings: []string{"a"}}, true},
		{"unknown kind", ast.Definition{Kind: ast.DefKind(99)}, true},
	}
	for _, c := range cases {
		if got := definitionFieldsInvalid(c.def); got != c.want {
			t.Errorf("%s: definitionFieldsInvalid = %v, want %v", c.name, got, c.want)
		}
	}
}
