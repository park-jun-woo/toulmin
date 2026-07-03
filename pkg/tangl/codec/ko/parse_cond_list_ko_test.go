//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseCondListKo — empty, single seed, 그리고/또는 chaining, missing prefix, term error
package ko

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/pkg/tangl/codec/scan"
)

func TestParseCondListKo(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		if _, err := parseCondListKo(nil, "t.md"); err == nil {
			t.Fatal("expected error for empty list")
		}
	})

	t.Run("SingleSeed", func(t *testing.T) {
		items := []scan.Item{{Text: "`a` 가 비어 있으면", No: 1}}
		expr, err := parseCondListKo(items, "t.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if _, ok := expr.(ast.Compare); !ok {
			t.Fatalf("got %#v, want Compare", expr)
		}
	})

	t.Run("AndChain", func(t *testing.T) {
		items := []scan.Item{
			{Text: "`a` 가 비어 있으면", No: 1},
			{Text: "그리고 `b` 가 비어 있으면", No: 2},
		}
		expr, err := parseCondListKo(items, "t.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		logic, ok := expr.(ast.Logic)
		if !ok || logic.Op != "and" || len(logic.Terms) != 2 {
			t.Fatalf("got %#v, want Logic{and,2}", expr)
		}
	})

	t.Run("OrChain", func(t *testing.T) {
		items := []scan.Item{
			{Text: "`a` 가 비어 있으면", No: 1},
			{Text: "또는 `직원` 이면", No: 2},
		}
		expr, err := parseCondListKo(items, "t.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		logic, ok := expr.(ast.Logic)
		if !ok || logic.Op != "or" {
			t.Fatalf("got %#v, want Logic{or}", expr)
		}
	})

	t.Run("MissingPrefix", func(t *testing.T) {
		items := []scan.Item{
			{Text: "`a` 가 비어 있으면", No: 1},
			{Text: "`b` 가 비어 있으면", No: 2},
		}
		_, err := parseCondListKo(items, "t.md")
		if err == nil || !strings.Contains(err.Error(), "그리고") {
			t.Fatalf("expected missing-prefix error, got %v", err)
		}
	})

	t.Run("SeedError", func(t *testing.T) {
		items := []scan.Item{{Text: "바나나", No: 1}}
		if _, err := parseCondListKo(items, "t.md"); err == nil {
			t.Fatal("expected error from a malformed seed term")
		}
	})

	t.Run("ChainedTermError", func(t *testing.T) {
		items := []scan.Item{
			{Text: "`a` 가 비어 있으면", No: 1},
			{Text: "그리고 바나나", No: 2},
		}
		if _, err := parseCondListKo(items, "t.md"); err == nil {
			t.Fatal("expected error from a malformed chained term")
		}
	})
}
