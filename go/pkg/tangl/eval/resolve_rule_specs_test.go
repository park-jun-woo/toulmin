//ff:what TestResolveRuleSpecs — optional DefinitionResolver materializes inferred rule Specs
package eval

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/go/pkg/toulmin"
)

type definitionSpecResolver struct{ fakeResolver }

type derivedSpec struct{ name string }

func (s *derivedSpec) SpecName() string { return s.name }
func (s *derivedSpec) Validate() error  { return nil }

func (definitionSpecResolver) DefinitionSpec(def ast.Definition) (toulmin.Spec, error) {
	return &derivedSpec{name: def.Name}, nil
}

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

func TestMergeRuleSpecsRejectsConflict(t *testing.T) {
	_, err := mergeRuleSpecs(toulmin.Specs{&derivedSpec{name: "x"}}, toulmin.Specs{&derivedSpec{name: "x"}})
	if err != nil {
		t.Fatalf("equal specs conflict: %v", err)
	}
	_, err = mergeRuleSpecs(toulmin.Specs{&derivedSpec{name: "x"}}, toulmin.Specs{&testSpec{Name: "x"}})
	if err == nil {
		t.Fatal("different specs with the same SpecName must conflict")
	}
}
