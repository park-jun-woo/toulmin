//ff:func feature=tangl type=validator control=sequence
//ff:what TestCheckDuplicateEndpointNames — tests checkDuplicateEndpointNames for empty, unique, and duplicate-name branches via subtests
package validate

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestCheckDuplicateEndpointNames(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		doc := &ast.Document{}
		errs := checkDuplicateEndpointNames(doc)
		if len(errs) != 0 {
			t.Fatalf("expected no errors, got %v", errs)
		}
	})

	t.Run("Unique", func(t *testing.T) {
		doc := &ast.Document{
			Provides: []ast.Endpoint{
				{Name: "x", Line: 1},
				{Name: "y", Line: 2},
			},
		}
		errs := checkDuplicateEndpointNames(doc)
		if len(errs) != 0 {
			t.Fatalf("expected no errors, got %v", errs)
		}
	})

	t.Run("Duplicate", func(t *testing.T) {
		doc := &ast.Document{
			Provides: []ast.Endpoint{
				{Name: "x", Line: 1},
				{Name: "x", Line: 5},
			},
		}
		errs := checkDuplicateEndpointNames(doc)
		if len(errs) != 1 {
			t.Fatalf("expected 1 error, got %v", errs)
		}
		if !strings.Contains(errs[0].Error(), "duplicate Endpoint name") {
			t.Errorf("expected duplicate Endpoint name error, got %v", errs[0])
		}
	})
}
