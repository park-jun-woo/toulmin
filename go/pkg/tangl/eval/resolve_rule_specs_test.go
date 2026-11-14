//ff:func feature=tangl type=engine control=sequence
//ff:what TestResolveRuleSpecs — optional DefinitionResolver materializes inferred rule Specs
package eval

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestResolveRuleSpecs(t *testing.T) {
	doc := &ast.Document{
		Defs:  []ast.Definition{{Name: "limit", Kind: ast.ConstDef, Value: "3", SpecRef: &ast.Ref{Name: "NewLimit"}}},
		Rules: []ast.InlineRule{{Name: "r", Cond: ast.Compare{Lhs: ast.SpecRef{Steps: []string{"limit", "Value"}}, Op: "equals", Value: "3"}}},
	}
	got, err := resolveRuleSpecs(doc, definitionSpecResolver{}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(got["r"]) != 1 || got["r"][0].SpecName() != "limit" {
		t.Fatalf("resolved=%+v", got)
	}
}
