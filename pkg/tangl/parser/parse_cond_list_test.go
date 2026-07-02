//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseCondList — tests parseCondList for empty-input, first-term-error, single-term success, missing-prefix, subsequent-term-error, and multi-term merge branches
package parser

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestParseCondList(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		_, err := parseCondList(nil, "test.md")
		if err == nil || !strings.Contains(err.Error(), "expected at least one condition") {
			t.Fatalf("expected empty-list error, got %v", err)
		}
	})

	t.Run("FirstTermError", func(t *testing.T) {
		items := []item{{Text: "bogus", Line: 1}}
		_, err := parseCondList(items, "test.md")
		if err == nil {
			t.Fatal("expected an error for a malformed first term")
		}
	})

	t.Run("SingleTermSuccess", func(t *testing.T) {
		items := []item{{Text: "`amount` is empty", Line: 1}}
		expr, err := parseCondList(items, "test.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		cmp, ok := expr.(ast.Compare)
		if !ok || cmp.Field != "amount" {
			t.Fatalf("expected Compare{amount}, got %+v", expr)
		}
	})

	t.Run("MissingPrefix", func(t *testing.T) {
		items := []item{
			{Text: "`amount` is empty", Line: 1},
			{Text: "`other` equals 5", Line: 2},
		}
		_, err := parseCondList(items, "test.md")
		if err == nil || !strings.Contains(err.Error(), "expected 'and'/'or' prefix") {
			t.Fatalf("expected missing-prefix error, got %v", err)
		}
	})

	t.Run("SubsequentTermError", func(t *testing.T) {
		items := []item{
			{Text: "`amount` is empty", Line: 1},
			{Text: "and bogus", Line: 2},
		}
		_, err := parseCondList(items, "test.md")
		if err == nil {
			t.Fatal("expected an error for a malformed subsequent term")
		}
	})

	t.Run("MultiTermSuccess", func(t *testing.T) {
		items := []item{
			{Text: "`amount` is empty", Line: 1},
			{Text: "and `other` equals 5", Line: 2},
		}
		expr, err := parseCondList(items, "test.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		logic, ok := expr.(ast.Logic)
		if !ok || logic.Op != "and" || len(logic.Terms) != 2 {
			t.Fatalf("expected Logic{and, 2 terms}, got %+v", expr)
		}
	})
}
