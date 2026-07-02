//ff:func feature=tangl type=validator control=sequence
//ff:what TestCaseNodeLineIndex — tests caseNodeLineIndex for empty-nodes and multiple-nodes branches
package validate

import (
	"reflect"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestCaseNodeLineIndex(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		c := ast.Case{Name: "x"}
		idx := caseNodeLineIndex(c)
		if len(idx) != 0 {
			t.Fatalf("expected empty index, got %v", idx)
		}
	})

	t.Run("Multiple", func(t *testing.T) {
		c := ast.Case{
			Name: "x",
			Nodes: []ast.Node{
				{Name: "a", Line: 3},
				{Name: "b", Line: 5},
			},
		}
		idx := caseNodeLineIndex(c)
		want := map[string]int{"a": 3, "b": 5}
		if !reflect.DeepEqual(idx, want) {
			t.Fatalf("expected idx=%v, got %v", want, idx)
		}
	})
}
