//ff:func feature=tangl type=validator control=sequence
//ff:what TestCheckDuplicateRuleNames — tests checkDuplicateRuleNames for empty, unique, and duplicate branches
package validate

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestCheckDuplicateRuleNames(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		errs := checkDuplicateRuleNames(&ast.Document{})
		if len(errs) != 0 {
			t.Fatalf("expected no errors, got %v", errs)
		}
	})

	t.Run("Unique", func(t *testing.T) {
		doc := &ast.Document{Rules: []ast.InlineRule{
			{Name: "a", Line: 1},
			{Name: "b", Line: 2},
		}}
		errs := checkDuplicateRuleNames(doc)
		if len(errs) != 0 {
			t.Fatalf("expected no errors, got %v", errs)
		}
	})

	t.Run("Duplicate", func(t *testing.T) {
		doc := &ast.Document{Rules: []ast.InlineRule{
			{Name: "a", Line: 1},
			{Name: "a", Line: 4},
		}}
		errs := checkDuplicateRuleNames(doc)
		if len(errs) != 1 {
			t.Fatalf("expected 1 error, got %v", errs)
		}
		if !strings.Contains(errs[0].Error(), "duplicate Rules entry") {
			t.Errorf("expected duplicate Rules entry error, got %v", errs[0])
		}
	})
}
