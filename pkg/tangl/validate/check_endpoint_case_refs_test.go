//ff:func feature=tangl type=validator control=sequence
//ff:what TestCheckEndpointCaseRefs — tests checkEndpointCaseRefs for valid, missing-run-case, and missing-check-case branches via subtests
package validate

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestCheckEndpointCaseRefs(t *testing.T) {
	t.Run("NoEndpoints", func(t *testing.T) {
		doc := &ast.Document{}
		errs := checkEndpointCaseRefs(doc)
		if len(errs) != 0 {
			t.Fatalf("expected no errors, got %v", errs)
		}
	})

	t.Run("Valid", func(t *testing.T) {
		doc := &ast.Document{
			Cases: []ast.Case{
				{Name: "x"},
				{Name: "y"},
			},
			Provides: []ast.Endpoint{
				{Name: "ep", Runs: []string{"x"}, Checks: []string{"y"}, Line: 1},
			},
		}
		errs := checkEndpointCaseRefs(doc)
		if len(errs) != 0 {
			t.Fatalf("expected no errors, got %v", errs)
		}
	})

	t.Run("MissingRunCase", func(t *testing.T) {
		doc := &ast.Document{
			Provides: []ast.Endpoint{
				{Name: "ep", Runs: []string{"missing"}, Line: 2},
			},
		}
		errs := checkEndpointCaseRefs(doc)
		if len(errs) != 1 {
			t.Fatalf("expected 1 error, got %v", errs)
		}
		if !strings.Contains(errs[0].Error(), "run case") {
			t.Errorf("expected run case error, got %v", errs[0])
		}
	})

	t.Run("MissingCheckCase", func(t *testing.T) {
		doc := &ast.Document{
			Provides: []ast.Endpoint{
				{Name: "ep", Checks: []string{"missing"}, Line: 3},
			},
		}
		errs := checkEndpointCaseRefs(doc)
		if len(errs) != 1 {
			t.Fatalf("expected 1 error, got %v", errs)
		}
		if !strings.Contains(errs[0].Error(), "check case") {
			t.Errorf("expected check case error, got %v", errs[0])
		}
	})
}
