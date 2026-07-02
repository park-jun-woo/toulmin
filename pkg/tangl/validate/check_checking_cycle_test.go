//ff:func feature=tangl type=validator control=sequence
//ff:what TestCheckCheckingCycle — tests checkCheckingCycle for no-cycle and cycle-detected branches via subtests
package validate

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestCheckCheckingCycle(t *testing.T) {
	t.Run("NoCycle", func(t *testing.T) {
		doc := &ast.Document{
			Cases: []ast.Case{
				{
					Name: "x",
					Line: 1,
					Nodes: []ast.Node{
						{Name: "a", Checking: "y"},
					},
				},
				{
					Name: "y",
					Line: 2,
				},
			},
		}
		errs := checkCheckingCycle(doc)
		if len(errs) != 0 {
			t.Fatalf("expected no errors, got %v", errs)
		}
	})

	t.Run("DetectsCycle", func(t *testing.T) {
		doc := &ast.Document{
			Cases: []ast.Case{
				{
					Name: "x",
					Line: 1,
					Nodes: []ast.Node{
						{Name: "a", Checking: "y"},
					},
				},
				{
					Name: "y",
					Line: 2,
					Nodes: []ast.Node{
						{Name: "b", Checking: "x"},
					},
				},
			},
		}
		errs := checkCheckingCycle(doc)
		if len(errs) != 1 {
			t.Fatalf("expected exactly 1 error, got %v", errs)
		}
		if !strings.Contains(errs[0].Error(), "checking cycle") {
			t.Errorf("expected checking cycle error, got %v", errs[0])
		}
	})
}
