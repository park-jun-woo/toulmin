//ff:func feature=tangl type=validator control=sequence
//ff:what TestCheckDuplicateCaseNames — tests checkDuplicateCaseNames for empty, unique, and duplicate-name branches via subtests
package validate

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestCheckDuplicateCaseNames(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		doc := &ast.Document{}
		errs := checkDuplicateCaseNames(doc)
		if len(errs) != 0 {
			t.Fatalf("expected no errors, got %v", errs)
		}
	})

	t.Run("Unique", func(t *testing.T) {
		doc := &ast.Document{
			Cases: []ast.Case{
				{Name: "x", Line: 1},
				{Name: "y", Line: 2},
			},
		}
		errs := checkDuplicateCaseNames(doc)
		if len(errs) != 0 {
			t.Fatalf("expected no errors, got %v", errs)
		}
	})

	t.Run("Duplicate", func(t *testing.T) {
		doc := &ast.Document{
			Cases: []ast.Case{
				{Name: "x", Line: 1},
				{Name: "x", Line: 5},
			},
		}
		errs := checkDuplicateCaseNames(doc)
		if len(errs) != 1 {
			t.Fatalf("expected 1 error, got %v", errs)
		}
		if !strings.Contains(errs[0].Error(), "duplicate case name") {
			t.Errorf("expected duplicate case name error, got %v", errs[0])
		}
	})
}
