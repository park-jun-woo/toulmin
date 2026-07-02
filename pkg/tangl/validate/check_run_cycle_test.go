//ff:func feature=tangl type=validator control=sequence
//ff:what TestCheckRunCycle — tests checkRunCycle for no-cycle and cycle-detected branches
package validate

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestCheckRunCycle(t *testing.T) {
	t.Run("NoCycle", func(t *testing.T) {
		doc := &ast.Document{
			Cases: []ast.Case{
				{
					Name: "x",
					Line: 1,
					Execs: []ast.Exec{
						{Kind: ast.RunExec, Case: "y"},
					},
				},
				{
					Name: "y",
					Line: 2,
				},
			},
		}
		errs := checkRunCycle(doc)
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
					Execs: []ast.Exec{
						{Kind: ast.RunExec, Case: "y"},
					},
				},
				{
					Name: "y",
					Line: 2,
					Execs: []ast.Exec{
						{Kind: ast.RunExec, Case: "x"},
					},
				},
			},
		}
		errs := checkRunCycle(doc)
		if len(errs) != 1 {
			t.Fatalf("expected exactly 1 error, got %v", errs)
		}
		if !strings.Contains(errs[0].Error(), "run cycle") {
			t.Errorf("expected run cycle error, got %v", errs[0])
		}
	})
}
