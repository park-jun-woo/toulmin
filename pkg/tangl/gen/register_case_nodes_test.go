//ff:func feature=tangl type=codegen control=sequence
//ff:what TestRegisterCaseNodes — tests registerCaseNodes for plain nodes, with-terms (spec/const/bare), qualified nodes, and role variants
package gen

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestRegisterCaseNodes(t *testing.T) {
	t.Run("plain rule node with no with/qualified", func(t *testing.T) {
		var w strings.Builder
		gc := &genContext{Defs: map[string]defInfo{}}
		c := ast.Case{Nodes: []ast.Node{{Name: "nodeOne", Role: ast.GeneralRule}}}
		nodes := registerCaseNodes(&w, gc, c)
		want := "\tnodeOne := g.Rule(nodeOne)\n"
		if w.String() != want {
			t.Errorf("got %q, want %q", w.String(), want)
		}
		if _, ok := nodes["nodeOne"]; !ok {
			t.Errorf("expected nodes map to contain nodeOne")
		}
	})

	t.Run("counter role node", func(t *testing.T) {
		var w strings.Builder
		gc := &genContext{Defs: map[string]defInfo{}}
		c := ast.Case{Nodes: []ast.Node{{Name: "n2", Role: ast.CounterRule}}}
		registerCaseNodes(&w, gc, c)
		if !strings.Contains(w.String(), "g.Counter(") {
			t.Errorf("expected Counter call, got %q", w.String())
		}
	})

	t.Run("except role node", func(t *testing.T) {
		var w strings.Builder
		gc := &genContext{Defs: map[string]defInfo{}}
		c := ast.Case{Nodes: []ast.Node{{Name: "n3", Role: ast.ExceptRule}}}
		registerCaseNodes(&w, gc, c)
		if !strings.Contains(w.String(), "g.Except(") {
			t.Errorf("expected Except call, got %q", w.String())
		}
	})

	t.Run("with terms spec, const, and bare identifier", func(t *testing.T) {
		var w strings.Builder
		gc := &genContext{Defs: map[string]defInfo{
			"specTerm":  {Spec: "specTermSpec"},
			"constTerm": {Const: "constTermConst"},
		}}
		c := ast.Case{Nodes: []ast.Node{
			{Name: "nodeWith", Role: ast.GeneralRule, With: []string{"specTerm", "constTerm", "bareTerm"}},
		}}
		registerCaseNodes(&w, gc, c)
		out := w.String()
		if !strings.Contains(out, ".With(specTermSpec)") {
			t.Errorf("missing spec with-arg: %q", out)
		}
		if !strings.Contains(out, ".With(constTermConst)") {
			t.Errorf("missing const with-arg: %q", out)
		}
		if !strings.Contains(out, ".With(bareTerm)") {
			t.Errorf("missing bare with-arg: %q", out)
		}
	})

	t.Run("qualified node emits qualifier call", func(t *testing.T) {
		var w strings.Builder
		gc := &genContext{Defs: map[string]defInfo{}}
		q := 0.5
		c := ast.Case{Nodes: []ast.Node{{Name: "nodeQ", Role: ast.GeneralRule, Qualified: &q}}}
		registerCaseNodes(&w, gc, c)
		if !strings.Contains(w.String(), ".Qualifier(0.5)") {
			t.Errorf("expected Qualifier call, got %q", w.String())
		}
	})

	t.Run("multiple nodes populate map for each", func(t *testing.T) {
		var w strings.Builder
		gc := &genContext{Defs: map[string]defInfo{}}
		c := ast.Case{Nodes: []ast.Node{
			{Name: "one", Role: ast.GeneralRule},
			{Name: "two", Role: ast.GeneralRule},
		}}
		nodes := registerCaseNodes(&w, gc, c)
		if len(nodes) != 2 {
			t.Fatalf("expected 2 nodes, got %d", len(nodes))
		}
	})
}
