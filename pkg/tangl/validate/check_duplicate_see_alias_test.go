//ff:func feature=tangl type=validator control=sequence
//ff:what TestCheckDuplicateSeeAlias — tests checkDuplicateSeeAlias for empty, unique, and duplicate-alias branches via subtests
package validate

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestCheckDuplicateSeeAlias(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		doc := &ast.Document{}
		errs := checkDuplicateSeeAlias(doc)
		if len(errs) != 0 {
			t.Fatalf("expected no errors, got %v", errs)
		}
	})

	t.Run("Unique", func(t *testing.T) {
		doc := &ast.Document{
			Sees: []ast.See{
				{Alias: "x", Line: 1},
				{Alias: "y", Line: 2},
			},
		}
		errs := checkDuplicateSeeAlias(doc)
		if len(errs) != 0 {
			t.Fatalf("expected no errors, got %v", errs)
		}
	})

	t.Run("Duplicate", func(t *testing.T) {
		doc := &ast.Document{
			Sees: []ast.See{
				{Alias: "x", Line: 1},
				{Alias: "x", Line: 5},
			},
		}
		errs := checkDuplicateSeeAlias(doc)
		if len(errs) != 1 {
			t.Fatalf("expected 1 error, got %v", errs)
		}
		if !strings.Contains(errs[0].Error(), "duplicate See alias") {
			t.Errorf("expected duplicate See alias error, got %v", errs[0])
		}
	})
}
