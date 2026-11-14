//ff:func feature=tangl type=codegen control=sequence
//ff:what TestRegisterCaseNodes_Grounded — tests registerCaseNodes Grounds emission for present and absent cases
package gen

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestRegisterCaseNodes_Grounded(t *testing.T) {
	t.Run("grounded node emits grounds call in document order", func(t *testing.T) {
		var w strings.Builder
		gc := &genContext{Defs: map[string]defInfo{}}
		c := ast.Case{Nodes: []ast.Node{
			{Name: "nodeG", Role: ast.CounterRule, Grounds: []string{"source-body", "mx"}},
		}}
		registerCaseNodes(&w, gc, c)
		if !strings.Contains(w.String(), ".Grounds(\"source-body\", \"mx\")") {
			t.Errorf("expected Grounds call, got %q", w.String())
		}
	})

	t.Run("node without grounds emits no grounds call", func(t *testing.T) {
		var w strings.Builder
		gc := &genContext{Defs: map[string]defInfo{}}
		c := ast.Case{Nodes: []ast.Node{{Name: "nodePlain", Role: ast.GeneralRule}}}
		registerCaseNodes(&w, gc, c)
		if strings.Contains(w.String(), ".Grounds(") {
			t.Errorf("expected no Grounds call, got %q", w.String())
		}
	})
}
