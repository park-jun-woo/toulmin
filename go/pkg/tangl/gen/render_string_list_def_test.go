//ff:func feature=tangl type=codegen control=sequence
//ff:what TestRenderStringListDef — a string list renders as a quoted NewStringListSpec var and reports its Spec ident
package gen

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestRenderStringListDef(t *testing.T) {
	var w strings.Builder
	info := renderStringListDef(&w, ast.Definition{Name: "allowed domains", Kind: ast.StringListDef}, []string{`a"b`, "c"})
	want := "var allowedDomainsSpec = toulmin.NewStringListSpec(\"allowed domains\", []string{\"a\\\"b\", \"c\"})\n\n"
	if w.String() != want {
		t.Fatalf("rendered:\n%q\nwant:\n%q", w.String(), want)
	}
	if info.Spec != "allowedDomainsSpec" || info.Kind != ast.StringListDef {
		t.Fatalf("defInfo = %+v", info)
	}
	w.Reset()
	renderStringListDef(&w, ast.Definition{Name: "none", Kind: ast.StringListDef}, nil)
	if !strings.Contains(w.String(), "[]string{})") {
		t.Fatalf("empty list rendered %q", w.String())
	}
}
