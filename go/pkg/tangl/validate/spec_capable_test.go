//ff:func feature=tangl type=validator control=iteration dimension=1
//ff:what TestSpecCapable — lists, files and constructor consts produce Specs; plain consts and structs do not
package validate

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestSpecCapable(t *testing.T) {
	for _, c := range []struct {
		def  ast.Definition
		want bool
	}{
		{ast.Definition{Kind: ast.StringListDef}, true},
		{ast.Definition{Kind: ast.FileSpecDef}, true},
		{ast.Definition{Kind: ast.ConstDef, SpecRef: &ast.Ref{Name: "NewX"}}, true},
		{ast.Definition{Kind: ast.ConstDef, Value: "1"}, false},
		{ast.Definition{Kind: ast.StructDef, SpecRef: &ast.Ref{Name: "NewX"}}, false},
	} {
		if got := specCapable(c.def); got != c.want {
			t.Errorf("specCapable(%+v) = %v, want %v", c.def, got, c.want)
		}
	}
}
