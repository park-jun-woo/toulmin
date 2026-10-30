//ff:func feature=tangl type=validator control=sequence
//ff:what TestCheckUsingAliasRefsRules — host rules must use a declared See alias; bare and declared refs pass
package validate

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestCheckUsingAliasRefsRules(t *testing.T) {
	doc := &ast.Document{
		Path: "p.md",
		Sees: []ast.See{{Alias: "ratchet", Package: "example/ratchet"}},
		Rules: []ast.InlineRule{
			{Name: "declared", Line: 2, Using: &ast.Ref{Alias: "ratchet", Name: "f"}},
			{Name: "bare", Line: 3, Using: &ast.Ref{Name: "local"}},
			{Name: "expr", Line: 4, Cond: ast.Compare{Field: "x", Op: "is empty"}},
			{Name: "stray", Line: 5, Using: &ast.Ref{Alias: "nlp", Name: "g"}},
		},
	}
	errs := checkUsingAliasRefs(doc)
	if len(errs) != 1 || errs[0].Error() != `p.md:5: rule "stray" uses undeclared package alias "nlp"` {
		t.Fatalf("errs = %v, want only the stray alias", errs)
	}
}
