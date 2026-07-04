//ff:func feature=tangl type=codegen control=sequence
//ff:what TestCollectCheckingTargets — tests collectCheckingTargets for empty, with-checking, and without-checking branches
package gen

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestCollectCheckingTargets(t *testing.T) {
	t.Run("every case is a target", func(t *testing.T) {
		doc := &ast.Document{
			Cases: []ast.Case{
				{
					Name: "case one",
					Nodes: []ast.Node{
						{Name: "n1", Checking: &ast.Ref{Name: "case two"}},
						{Name: "n2"},
					},
				},
				{
					Name:  "case two",
					Nodes: []ast.Node{},
				},
			},
		}
		targets := collectCheckingTargets(doc)
		if len(targets) != 2 || !targets["case one"] || !targets["case two"] {
			t.Errorf("expected every case name, got %+v", targets)
		}
	})

	t.Run("empty", func(t *testing.T) {
		targets := collectCheckingTargets(&ast.Document{})
		if len(targets) != 0 {
			t.Errorf("expected empty targets, got %+v", targets)
		}
	})
}
