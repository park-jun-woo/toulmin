package eval

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestCompileDocumentStringListMembership(t *testing.T) {
	doc := &ast.Document{
		Defs:  []ast.Definition{{Name: "domains", Kind: ast.StringListDef, Strings: []string{"a,b", "c"}}},
		Rules: []ast.InlineRule{{Name: "listed", Cond: ast.Compare{Lhs: ast.Path{Steps: []string{"domain"}}, Op: "is in", Rhs: ast.SpecRef{Steps: []string{"domains"}}}}},
	}
	rules, err := CompileDoc(doc, nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx := toulmin.NewContext()
	ctx.Set("domain", "a,b")
	got, _ := rules["listed"](ctx, nil)
	if !got {
		t.Fatal("listed value was rejected")
	}
	ctx.Set("domain", "a")
	got, _ = rules["listed"](ctx, nil)
	if got {
		t.Fatal("partial value was accepted")
	}
}
