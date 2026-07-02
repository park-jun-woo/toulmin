//ff:func feature=tangl type=validator control=sequence
//ff:what TestCaseNodeSet — tests caseNodeSet for empty-nodes and multiple-nodes branches
package validate

import (
	"reflect"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestCaseNodeSet(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		c := ast.Case{Name: "x"}
		set := caseNodeSet(c)
		if len(set) != 0 {
			t.Fatalf("expected empty set, got %v", set)
		}
	})

	t.Run("Multiple", func(t *testing.T) {
		c := ast.Case{
			Name: "x",
			Nodes: []ast.Node{
				{Name: "a"},
				{Name: "b"},
			},
		}
		set := caseNodeSet(c)
		want := map[string]bool{"a": true, "b": true}
		if !reflect.DeepEqual(set, want) {
			t.Fatalf("expected set=%v, got %v", want, set)
		}
	})
}
