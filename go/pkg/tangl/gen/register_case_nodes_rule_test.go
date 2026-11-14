//ff:func feature=tangl type=codegen control=sequence
//ff:what TestRegisterCaseNodesRuleBinding — a node bound to a rule inherits its With, inferred Specs and Grounds once
package gen

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestRegisterCaseNodesRuleBinding(t *testing.T) {
	rule := ast.InlineRule{
		Name:    "fits",
		With:    []string{"limit"},
		Grounds: []string{"usage"},
		Cond:    ast.Compare{Lhs: ast.SpecRef{Steps: []string{"cap"}}, Op: "is not empty"},
	}
	gc := &genContext{Rules: map[string]ast.InlineRule{"fits": rule}}
	var w strings.Builder
	nodes, err := registerCaseNodes(&w, gc, ast.Case{Name: "c", Nodes: []ast.Node{{Name: "fits", Role: ast.GeneralRule}}})
	if err != nil {
		t.Fatal(err)
	}
	line := w.String()
	if !strings.Contains(line, ".With(limit).With(cap)") || !strings.Contains(line, `.Grounds(`) || !strings.Contains(line, `"usage"`) {
		t.Fatalf("inherited bindings missing: %q", line)
	}
	if nodes["fits"].Var != "fits" {
		t.Fatalf("nodes = %+v", nodes)
	}
}
