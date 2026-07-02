//ff:func feature=tangl type=validator control=sequence
//ff:what TestCheckExecNodeRefs — tests checkExecNodeRefs for no-cases, valid node ref, and missing node ref branches
package validate

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestCheckExecNodeRefs(t *testing.T) {
	t.Run("NoCases", func(t *testing.T) {
		doc := &ast.Document{}
		errs := checkExecNodeRefs(doc)
		if len(errs) != 0 {
			t.Fatalf("expected no errors, got %v", errs)
		}
	})

	t.Run("Valid", func(t *testing.T) {
		doc := &ast.Document{
			Cases: []ast.Case{
				{
					Name:  "x",
					Nodes: []ast.Node{{Name: "a"}},
					Execs: []ast.Exec{
						{Node: "a", Line: 1},
					},
				},
			},
		}
		errs := checkExecNodeRefs(doc)
		if len(errs) != 0 {
			t.Fatalf("expected no errors, got %v", errs)
		}
	})

	t.Run("Missing", func(t *testing.T) {
		doc := &ast.Document{
			Cases: []ast.Case{
				{
					Name: "x",
					Execs: []ast.Exec{
						{Node: "missing", Line: 4},
					},
				},
			},
		}
		errs := checkExecNodeRefs(doc)
		if len(errs) != 1 {
			t.Fatalf("expected 1 error, got %v", errs)
		}
		if !strings.Contains(errs[0].Error(), "is not a registered node") {
			t.Errorf("expected not-registered-node error, got %v", errs[0])
		}
	})
}
