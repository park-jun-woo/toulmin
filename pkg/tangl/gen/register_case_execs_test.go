//ff:func feature=tangl type=codegen control=sequence
//ff:what TestRegisterCaseExecs — tests registerCaseExecs for node-without-execs skip, successful render, and error propagation branches
package gen

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestRegisterCaseExecs(t *testing.T) {
	t.Run("node without matching execs is skipped", func(t *testing.T) {
		var w strings.Builder
		nodes := map[string]nodeInfo{
			"a": {Var: "nA", Node: ast.Node{Name: "a"}},
		}
		c := ast.Case{
			Name:  "caseA",
			Nodes: []ast.Node{{Name: "a"}},
		}
		err := registerCaseExecs(&w, "subj", c, nodes)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if w.String() != "\n" {
			t.Errorf("got %q, want %q", w.String(), "\n")
		}
	})

	t.Run("node with matching exec renders block", func(t *testing.T) {
		var w strings.Builder
		nodes := map[string]nodeInfo{
			"a": {Var: "nA", Node: ast.Node{Name: "a"}},
		}
		c := ast.Case{
			Name:  "caseA",
			Nodes: []ast.Node{{Name: "a"}},
			Execs: []ast.Exec{
				{Node: "a", Kind: ast.RunExec, Case: "otherCase"},
			},
		}
		err := registerCaseExecs(&w, "subj", c, nodes)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !strings.Contains(w.String(), "nA") || !strings.Contains(w.String(), "otherCaseGraph") {
			t.Errorf("unexpected output: %q", w.String())
		}
	})

	t.Run("error from runTarget propagates", func(t *testing.T) {
		var w strings.Builder
		nodes := map[string]nodeInfo{
			"a": {Var: "nA", Node: ast.Node{Name: "a"}},
		}
		c := ast.Case{
			Name:  "caseA",
			Nodes: []ast.Node{{Name: "a"}},
			Execs: []ast.Exec{
				{Node: "a", Kind: ast.RunExec, Case: "caseOne"},
				{Node: "a", Kind: ast.RunExec, Case: "caseTwo"},
			},
		}
		err := registerCaseExecs(&w, "subj", c, nodes)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}
