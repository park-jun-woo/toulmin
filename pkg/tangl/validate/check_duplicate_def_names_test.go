//ff:func feature=tangl type=validator control=sequence
//ff:what TestCheckDuplicateDefNames — tests checkDuplicateDefNames for empty, unique, and duplicate-name branches via subtests
package validate

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestCheckDuplicateDefNames(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		doc := &ast.Document{}
		errs := checkDuplicateDefNames(doc)
		if len(errs) != 0 {
			t.Fatalf("expected no errors, got %v", errs)
		}
	})

	t.Run("Unique", func(t *testing.T) {
		doc := &ast.Document{
			Defs: []ast.Definition{
				{Name: "x", Line: 1},
				{Name: "y", Line: 2},
			},
		}
		errs := checkDuplicateDefNames(doc)
		if len(errs) != 0 {
			t.Fatalf("expected no errors, got %v", errs)
		}
	})

	t.Run("Duplicate", func(t *testing.T) {
		doc := &ast.Document{
			Defs: []ast.Definition{
				{Name: "x", Line: 1},
				{Name: "x", Line: 5},
			},
		}
		errs := checkDuplicateDefNames(doc)
		if len(errs) != 1 {
			t.Fatalf("expected 1 error, got %v", errs)
		}
		if !strings.Contains(errs[0].Error(), "duplicate Definitions term") {
			t.Errorf("expected duplicate Definitions term error, got %v", errs[0])
		}
	})
}
