//ff:func feature=tangl type=codegen control=sequence
//ff:what TestWriteUnusedNodeGuards — tests writeUnusedNodeGuards for used-node-skip and unused-node-guard branches
package gen

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestWriteUnusedNodeGuards(t *testing.T) {
	t.Run("all used", func(t *testing.T) {
		var w strings.Builder
		c := ast.Case{
			Nodes: []ast.Node{
				{Name: "nodeA"},
			},
			Attacks: []ast.Attack{
				{Target: "nodeA", Attacker: "nodeA"},
			},
		}
		nodes := map[string]nodeInfo{
			"nodeA": {Var: "nodeAVar"},
		}
		writeUnusedNodeGuards(&w, c, nodes)
		if w.String() != "" {
			t.Errorf("expected empty output for used node, got %q", w.String())
		}
	})

	t.Run("unused node", func(t *testing.T) {
		var w strings.Builder
		c := ast.Case{
			Nodes: []ast.Node{
				{Name: "nodeA"},
				{Name: "nodeB"},
			},
			Execs: []ast.Exec{
				{Node: "nodeA"},
			},
		}
		nodes := map[string]nodeInfo{
			"nodeA": {Var: "nodeAVar"},
			"nodeB": {Var: "nodeBVar"},
		}
		writeUnusedNodeGuards(&w, c, nodes)
		got := w.String()
		if strings.Contains(got, "nodeAVar") {
			t.Errorf("expected nodeA to be skipped as used, got %q", got)
		}
		if !strings.Contains(got, "\t_ = nodeBVar\n") {
			t.Errorf("expected guard for unused nodeB, got %q", got)
		}
	})
}
