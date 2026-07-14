//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseCompareAdnKo — 관형형 leaf parse: table operators, bare 인/이 아닌, operand sides
package ko

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestParseCompareAdnKo(t *testing.T) {
	t.Run("PositiveOperator", func(t *testing.T) {
		expr, err := parseCompareAdnKo("`c` 의 `최소가격` 이 `기본가` 이하인", 1, "t.md")
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		cmp := expr.(ast.Compare)
		p := cmp.Lhs.(ast.Path)
		if len(p.Steps) != 2 || p.Steps[1] != "최소가격" || cmp.Op != "is at most" {
			t.Fatalf("got %#v", cmp)
		}
		if _, ok := cmp.Rhs.(ast.Path); !ok {
			t.Fatalf("Rhs = %#v", cmp.Rhs)
		}
	})

	t.Run("QuotedLiteral", func(t *testing.T) {
		expr, err := parseCompareAdnKo("`f` 의 `이름` 이 \"main.go\" 과 같은", 1, "t.md")
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		cmp := expr.(ast.Compare)
		if cmp.Value != "main.go" {
			t.Fatalf("Value = %q, want %q", cmp.Value, "main.go")
		}
	})

	t.Run("NegatedOperator", func(t *testing.T) {
		expr, err := parseCompareAdnKo("`점수` 가 90 이하가 아닌", 1, "t.md")
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		not := expr.(ast.Not)
		if cmp := not.Term.(ast.Compare); cmp.Op != "is at most" || cmp.Field != "점수" {
			t.Fatalf("inner = %#v", not.Term)
		}
	})

	t.Run("BareBooleanPositive", func(t *testing.T) {
		expr, err := parseCompareAdnKo("`활성` 인", 1, "t.md")
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		cmp := expr.(ast.Compare)
		if cmp.Field != "활성" || cmp.Op != "equals" || cmp.Value != "true" {
			t.Fatalf("got %#v", cmp)
		}
	})

	t.Run("BareBooleanNegated", func(t *testing.T) {
		expr, err := parseCompareAdnKo("`활성` 이 아닌", 1, "t.md")
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		not := expr.(ast.Not)
		if cmp := not.Term.(ast.Compare); cmp.Field != "활성" || cmp.Op != "equals" {
			t.Fatalf("inner = %#v", not.Term)
		}
	})

	t.Run("BarePathBoolean", func(t *testing.T) {
		expr, err := parseCompareAdnKo("`작성자` 의 `활성` 인", 1, "t.md")
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		cmp := expr.(ast.Compare)
		if _, ok := cmp.Lhs.(ast.Path); !ok || cmp.Op != "equals" {
			t.Fatalf("got %#v", cmp)
		}
	})

	t.Run("MalformedLhsOperand", func(t *testing.T) {
		if _, err := parseCompareAdnKo("기준 이상인", 1, "t.md"); err == nil {
			t.Fatal("expected error: 기준 without a backtick-quoted spec term")
		}
	})

	t.Run("MissingLhsField", func(t *testing.T) {
		if _, err := parseCompareAdnKo("점수 이상인", 1, "t.md"); err == nil {
			t.Fatal("expected error: no backtick-quoted field before the operator")
		}
	})

	t.Run("MalformedRhsOperand", func(t *testing.T) {
		if _, err := parseCompareAdnKo("`a` 가 기준 와 같은", 1, "t.md"); err == nil {
			t.Fatal("expected error: 기준 right side without a backtick-quoted spec term")
		}
	})

	t.Run("Unrecognized", func(t *testing.T) {
		if _, err := parseCompareAdnKo("`x` 가 이상함", 1, "t.md"); err == nil {
			t.Fatal("expected error for an unrecognized adnominal comparison")
		}
	})
}
