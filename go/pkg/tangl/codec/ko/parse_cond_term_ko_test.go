//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseCondTermKo — 다음 중 group, bare boolean name, 이 아니고 negation, compare, unrecognized
package ko

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/codec/scan"
)

func TestParseCondTermKo(t *testing.T) {
	t.Run("EitherGroup", func(t *testing.T) {
		children := []scan.Item{
			{Text: "`a` 가 비어 있으면", No: 2},
			{Text: "또는 `b` 가 비어 있으면", No: 3},
		}
		expr, err := parseCondTermKo("다음 중", children, 1, "t.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		either, ok := expr.(ast.Either)
		if !ok || len(either.Terms) != 1 {
			t.Fatalf("got %#v, want Either", expr)
		}
		if _, ok := either.Terms[0].(ast.Logic); !ok {
			t.Fatalf("either inner = %#v, want Logic", either.Terms[0])
		}
	})

	t.Run("BareBooleanName", func(t *testing.T) {
		expr, err := parseCondTermKo("`직원` 이면", nil, 1, "t.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		want := ast.Compare{Field: "직원", Op: "equals", Value: "true"}
		if expr != want {
			t.Fatalf("got %#v, want %#v", expr, want)
		}
	})

	t.Run("BareNameNegation", func(t *testing.T) {
		expr, err := parseCondTermKo("`직원` 이 아니고", nil, 1, "t.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		not, ok := expr.(ast.Not)
		if !ok {
			t.Fatalf("got %#v, want Not", expr)
		}
		if cmp, _ := not.Term.(ast.Compare); cmp.Field != "직원" || cmp.Op != "equals" || cmp.Value != "true" {
			t.Fatalf("inner = %#v", not.Term)
		}
	})

	t.Run("Comparison", func(t *testing.T) {
		expr, err := parseCondTermKo("`신용점수` 가 650 이상이면", nil, 1, "t.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		cmp, ok := expr.(ast.Compare)
		if !ok || cmp.Field != "신용점수" || cmp.Op != "is at least" || cmp.Value != "650" {
			t.Fatalf("got %#v", expr)
		}
	})

	t.Run("Unrecognized", func(t *testing.T) {
		if _, err := parseCondTermKo("바나나", nil, 1, "t.md"); err == nil {
			t.Fatal("expected error for an unrecognized condition")
		}
	})

	t.Run("EitherKeywordWithTrailingText", func(t *testing.T) {
		if _, err := parseCondTermKo("다음 중 골라", nil, 1, "t.md"); err == nil {
			t.Fatal("expected error: '다음 중' keyword with non-empty trailer is not an Either group")
		}
	})

	t.Run("EitherGroupChildError", func(t *testing.T) {
		if _, err := parseCondTermKo("다음 중", nil, 1, "t.md"); err == nil {
			t.Fatal("expected error propagated from an empty child list")
		}
	})

	t.Run("NegationWithoutBacktick", func(t *testing.T) {
		if _, err := parseCondTermKo("직원 이 아니고", nil, 1, "t.md"); err == nil {
			t.Fatal("expected error: '아니고' suffix without a backtick-quoted field")
		}
	})

	t.Run("NegationWithTrailingText", func(t *testing.T) {
		if _, err := parseCondTermKo("`직원` 뭔가 아니고", nil, 1, "t.md"); err == nil {
			t.Fatal("expected error: text after the backtick name is not only particles")
		}
	})

	t.Run("BareNameWithTrailingText", func(t *testing.T) {
		if _, err := parseCondTermKo("`직원` 뭔가", nil, 1, "t.md"); err == nil {
			t.Fatal("expected error: text after the backtick name is not only particles")
		}
	})
}
