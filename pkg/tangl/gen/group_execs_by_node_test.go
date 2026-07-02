//ff:func feature=tangl type=codegen control=sequence
//ff:what TestGroupExecsByNode — tests groupExecsByNode for empty input, single node, and multi-node grouping
package gen

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestGroupExecsByNode(t *testing.T) {
	t.Run("empty slice yields empty map", func(t *testing.T) {
		got := groupExecsByNode(nil)
		if len(got) != 0 {
			t.Fatalf("expected empty map, got %v", got)
		}
	})

	t.Run("groups execs by node key", func(t *testing.T) {
		execs := []ast.Exec{
			{Node: "a", Line: 1},
			{Node: "b", Line: 2},
			{Node: "a", Line: 3},
		}
		got := groupExecsByNode(execs)
		if len(got) != 2 {
			t.Fatalf("expected 2 groups, got %d", len(got))
		}
		if len(got["a"]) != 2 {
			t.Errorf("expected 2 execs for node a, got %d", len(got["a"]))
		}
		if len(got["b"]) != 1 {
			t.Errorf("expected 1 exec for node b, got %d", len(got["b"]))
		}
		if got["a"][0].Line != 1 || got["a"][1].Line != 3 {
			t.Errorf("unexpected order for node a: %+v", got["a"])
		}
	})
}
