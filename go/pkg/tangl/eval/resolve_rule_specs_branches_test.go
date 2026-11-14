//ff:func feature=tangl type=engine control=sequence
//ff:what TestResolveRuleSpecsBranches — list terms win, non-constructor terms skip, loader errors and nil Specs fail
package eval

import (
	"errors"
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/go/pkg/toulmin"
)

func TestResolveRuleSpecsBranches(t *testing.T) {
	ctorDef := ast.Definition{Name: "limit", Kind: ast.ConstDef, Value: "3", SpecRef: &ast.Ref{Name: "NewLimit"}}
	plainDef := ast.Definition{Name: "plain", Kind: ast.ConstDef, Value: "1"}
	rule := ast.InlineRule{Name: "r", With: []string{"domains", "plain", "missing"}, Cond: ast.Compare{Lhs: ast.SpecRef{Steps: []string{"limit"}}, Op: "is empty"}}
	doc := &ast.Document{Defs: []ast.Definition{ctorDef, plainDef}, Rules: []ast.InlineRule{rule}}
	lists := map[string]toulmin.StringListSpec{"domains": toulmin.NewStringListSpec("domains", []string{"a"})}

	got, err := resolveRuleSpecs(doc, definitionSpecResolver{}, lists)
	if err != nil {
		t.Fatal(err)
	}
	if len(got["r"]) != 2 || got["r"][0].SpecName() != "domains" || got["r"][1].SpecName() != "limit" {
		t.Fatalf("resolved = %+v, want [domains limit]", got["r"])
	}

	got, err = resolveRuleSpecs(doc, fakeResolver{}, lists)
	if err != nil || len(got["r"]) != 1 || got["r"][0].SpecName() != "domains" {
		t.Fatalf("non-DefinitionResolver: got %+v err %v, want only the list", got["r"], err)
	}

	_, err = resolveRuleSpecs(doc, failingDefinitionResolver{err: errors.New("ctor broke")}, lists)
	if err == nil || !strings.Contains(err.Error(), `rule "r" definition "limit": ctor broke`) {
		t.Fatalf("loader error: %v", err)
	}
	_, err = resolveRuleSpecs(doc, failingDefinitionResolver{}, lists)
	if err == nil || !strings.Contains(err.Error(), "produced a nil Spec") {
		t.Fatalf("nil Spec: %v", err)
	}
}
