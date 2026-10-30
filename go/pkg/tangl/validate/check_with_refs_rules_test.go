//ff:func feature=tangl type=validator control=iteration dimension=1
//ff:what TestCheckWithRefsRules — rule With/SpecRef terms must be defined (when explicit) and Spec-capable; duplicates report once
package validate

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestCheckWithRefsRules(t *testing.T) {
	doc := &ast.Document{
		Path: "p.md",
		Defs: []ast.Definition{
			{Name: "domains", Kind: ast.StringListDef, Strings: []string{"a"}},
			{Name: "shape", Kind: ast.StructDef},
			{Name: "plain", Kind: ast.ConstDef, Value: "1"},
		},
		Rules: []ast.InlineRule{
			{Name: "ok", Line: 2, With: []string{"domains"}, Cond: ast.Compare{Lhs: ast.SpecRef{Steps: []string{"domains"}}, Op: "is not empty"}},
			{Name: "missing", Line: 3, With: []string{"ghost", "ghost"}, Cond: ast.Compare{Field: "x", Op: "is empty"}},
			{Name: "inferred undefined", Line: 4, Cond: ast.Compare{Lhs: ast.SpecRef{Steps: []string{"nowhere"}}, Op: "is empty"}},
			{Name: "not spec", Line: 5, With: []string{"shape"}, Cond: ast.Compare{Lhs: ast.SpecRef{Steps: []string{"plain"}}, Op: "is empty"}},
		},
	}
	errs := checkWithRefs(doc)
	want := []string{
		`p.md:3: rule "missing" references undefined definition "ghost"`,
		`p.md:5: rule "not spec" references definition "shape" that cannot produce a Spec`,
		`p.md:5: rule "not spec" references definition "plain" that cannot produce a Spec`,
	}
	if len(errs) != len(want) {
		t.Fatalf("errs = %v, want %d errors", errs, len(want))
	}
	for i, w := range want {
		if !strings.Contains(errs[i].Error(), w) {
			t.Errorf("errs[%d] = %v, want %q", i, errs[i], w)
		}
	}
}
