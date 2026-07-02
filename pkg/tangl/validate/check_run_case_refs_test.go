//ff:func feature=tangl type=validator control=sequence
//ff:what TestCheckRunCaseRefs — tests checkRunCaseRefs for no-cases, non-run-exec skip, valid, and missing-target branches
package validate

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestCheckRunCaseRefs(t *testing.T) {
	t.Run("NoCases", func(t *testing.T) {
		doc := &ast.Document{}
		errs := checkRunCaseRefs(doc)
		if len(errs) != 0 {
			t.Fatalf("expected no errors, got %v", errs)
		}
	})

	t.Run("SkipsNonRunExec", func(t *testing.T) {
		doc := &ast.Document{
			Cases: []ast.Case{
				{Name: "a", Execs: []ast.Exec{
					{Kind: ast.DoExec, Case: "missing", Line: 1},
				}},
			},
		}
		errs := checkRunCaseRefs(doc)
		if len(errs) != 0 {
			t.Fatalf("expected no errors for non-run exec, got %v", errs)
		}
	})

	t.Run("ValidRunExec", func(t *testing.T) {
		doc := &ast.Document{
			Cases: []ast.Case{
				{Name: "a"},
				{Name: "b", Execs: []ast.Exec{
					{Kind: ast.RunExec, Case: "a", Line: 2},
				}},
			},
		}
		errs := checkRunCaseRefs(doc)
		if len(errs) != 0 {
			t.Fatalf("expected no errors, got %v", errs)
		}
	})

	t.Run("MissingRunTarget", func(t *testing.T) {
		doc := &ast.Document{
			Cases: []ast.Case{
				{Name: "b", Execs: []ast.Exec{
					{Kind: ast.RunExec, Case: "missing", Line: 3},
				}},
			},
		}
		errs := checkRunCaseRefs(doc)
		if len(errs) != 1 {
			t.Fatalf("expected 1 error, got %v", errs)
		}
		if !strings.Contains(errs[0].Error(), "run target case") {
			t.Errorf("expected run target case error, got %v", errs[0])
		}
	})
}
