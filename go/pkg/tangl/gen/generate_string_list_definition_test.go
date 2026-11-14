//ff:func feature=tangl type=codegen control=sequence
//ff:what TestGenerateStringListDefinition — an inline string-list definition is emitted as NewStringListSpec
package gen

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

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
