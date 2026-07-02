//ff:func feature=tangl type=validator control=sequence
//ff:what TestCheckInternalCaseRefs — tests checkInternalCaseRefs for no-internals, valid, empty-until, missing-until, missing-run, and missing-check branches
package validate

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestCheckInternalCaseRefs(t *testing.T) {
	t.Run("NoInternals", func(t *testing.T) {
		doc := &ast.Document{}
		errs := checkInternalCaseRefs(doc)
		if len(errs) != 0 {
			t.Fatalf("expected no errors, got %v", errs)
		}
	})

	t.Run("Valid", func(t *testing.T) {
		doc := &ast.Document{
			Cases: []ast.Case{
				{Name: "a"},
				{Name: "b"},
				{Name: "c"},
			},
			Internals: []ast.Internal{
				{Until: "a", Runs: []string{"b"}, Checks: []string{"c"}, Line: 1},
			},
		}
		errs := checkInternalCaseRefs(doc)
		if len(errs) != 0 {
			t.Fatalf("expected no errors, got %v", errs)
		}
	})

	t.Run("EmptyUntil", func(t *testing.T) {
		doc := &ast.Document{
			Internals: []ast.Internal{
				{Until: "", Line: 1},
			},
		}
		errs := checkInternalCaseRefs(doc)
		if len(errs) != 0 {
			t.Fatalf("expected no errors, got %v", errs)
		}
	})

	t.Run("MissingUntil", func(t *testing.T) {
		doc := &ast.Document{
			Internals: []ast.Internal{
				{Until: "missing", Line: 2},
			},
		}
		errs := checkInternalCaseRefs(doc)
		if len(errs) != 1 {
			t.Fatalf("expected 1 error, got %v", errs)
		}
		if !strings.Contains(errs[0].Error(), "until case") {
			t.Errorf("expected until case error, got %v", errs[0])
		}
	})

	t.Run("MissingRun", func(t *testing.T) {
		doc := &ast.Document{
			Internals: []ast.Internal{
				{Runs: []string{"missing"}, Line: 3},
			},
		}
		errs := checkInternalCaseRefs(doc)
		if len(errs) != 1 {
			t.Fatalf("expected 1 error, got %v", errs)
		}
		if !strings.Contains(errs[0].Error(), "run case") {
			t.Errorf("expected run case error, got %v", errs[0])
		}
	})

	t.Run("MissingCheck", func(t *testing.T) {
		doc := &ast.Document{
			Internals: []ast.Internal{
				{Checks: []string{"missing"}, Line: 4},
			},
		}
		errs := checkInternalCaseRefs(doc)
		if len(errs) != 1 {
			t.Fatalf("expected 1 error, got %v", errs)
		}
		if !strings.Contains(errs[0].Error(), "check case") {
			t.Errorf("expected check case error, got %v", errs[0])
		}
	})
}
