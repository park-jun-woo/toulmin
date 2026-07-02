//ff:func feature=tangl type=codegen control=sequence
//ff:what TestCollectCheckingTargets — tests collectCheckingTargets for empty, with-checking, and without-checking branches
package gen

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestCollectCheckingTargets(t *testing.T) {
	t.Run("with and without checking", func(t *testing.T) {
		doc := &ast.Document{
			Cases: []ast.Case{
				{
					Nodes: []ast.Node{
						{Name: "n1", Checking: "target one"},
						{Name: "n2"},
					},
				},
				{
					Nodes: []ast.Node{},
				},
			},
		}
		targets := collectCheckingTargets(doc)
		if len(targets) != 1 || !targets["target one"] {
			t.Errorf("expected only %q, got %+v", "target one", targets)
		}
	})

	t.Run("empty", func(t *testing.T) {
		targets := collectCheckingTargets(&ast.Document{})
		if len(targets) != 0 {
			t.Errorf("expected empty targets, got %+v", targets)
		}
	})
}
