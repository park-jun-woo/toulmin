//ff:func feature=tangl type=codegen control=sequence dimension=1
//ff:what TestRenderCase — tests renderCase for successful rendering and error propagation from registerCaseExecs
package gen

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestRenderCase(t *testing.T) {
	t.Run("successful render", func(t *testing.T) {
		var w strings.Builder
		gc := &genContext{
			Doc:  &ast.Document{Subject: "orders"},
			Defs: map[string]defInfo{},
		}
		c := ast.Case{
			Name:  "checkOrder",
			Nodes: []ast.Node{{Name: "nodeA", Role: ast.GeneralRule}},
		}
		err := renderCase(&w, gc, c)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		out := w.String()
		if !strings.Contains(out, "func newCheckOrderGraph() *toulmin.Graph {") {
			t.Errorf("missing function signature: %q", out)
		}
		if !strings.Contains(out, "var checkOrderGraph = newCheckOrderGraph()") {
			t.Errorf("missing var declaration: %q", out)
		}
	})

	t.Run("error from registerCaseExecs propagates", func(t *testing.T) {
		var w strings.Builder
		gc := &genContext{
			Doc:  &ast.Document{Subject: "orders"},
			Defs: map[string]defInfo{},
		}
		c := ast.Case{
			Name:  "checkOrder",
			Nodes: []ast.Node{{Name: "nodeA", Role: ast.GeneralRule}},
			Execs: []ast.Exec{
				{Node: "nodeA", Kind: ast.RunExec, Case: "caseOne"},
				{Node: "nodeA", Kind: ast.RunExec, Case: "caseTwo"},
			},
		}
		err := renderCase(&w, gc, c)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}
