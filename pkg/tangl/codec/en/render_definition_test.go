//ff:func feature=tangl type=parser control=sequence
//ff:what TestRenderDefinition — renderDefinition covers const (with/without as) and struct head
package en

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// TestRenderDefinition checks the constant forms (bare value and value-plus-as)
// and the struct definition head.
func TestRenderDefinition(t *testing.T) {
	bare := renderDefinition(ast.Definition{Name: "n", Kind: ast.ConstDef, Value: "650"})
	if bare != "`n` means 650" {
		t.Errorf("bare const = %q", bare)
	}
	withAs := renderDefinition(ast.Definition{
		Name: "credit threshold", Kind: ast.ConstDef, Value: "650",
		SpecRef: &ast.Ref{Alias: "credit", Name: "Threshold"},
	})
	if withAs != "`credit threshold` means 650 as `credit`.`Threshold`" {
		t.Errorf("const with as = %q", withAs)
	}
	strukt := renderDefinition(ast.Definition{Name: "applicant", Kind: ast.StructDef})
	if strukt != "`applicant` means" {
		t.Errorf("struct head = %q", strukt)
	}
	if got := renderDefinition(ast.Definition{Name: "domains", Kind: ast.StringListDef, Strings: []string{"a,b", "c"}}); got != "`domains` means [\"a,b\",\"c\"]" {
		t.Errorf("string list: got %q", got)
	}
	if got := renderDefinition(ast.Definition{Name: "domains", Kind: ast.FileSpecDef, LoadType: "JSON list", FilePath: "./x.json"}); got != "`domains` reads `JSON list` from \"./x.json\"" {
		t.Errorf("file spec: got %q", got)
	}
}
