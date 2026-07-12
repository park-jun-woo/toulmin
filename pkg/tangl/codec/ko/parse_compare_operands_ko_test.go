//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseCompareKoOperands — operand Lhs/Rhs forms (path/spec/call) and the Breaking Change (2) edges
package ko

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestParseCompareKoOperands(t *testing.T) {
	t.Run("PathLhsLiteralRhs", func(t *testing.T) {
		expr, err := parseCompareKo("`작성자` 의 `신뢰점수` 가 0.7 이상이면", 1, "t.md")
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		cmp := expr.(ast.Compare)
		p, ok := cmp.Lhs.(ast.Path)
		if !ok || len(p.Steps) != 2 || p.Steps[1] != "신뢰점수" || cmp.Field != "" {
			t.Fatalf("Lhs = %#v Field=%q", cmp.Lhs, cmp.Field)
		}
		if cmp.Op != "is at least" || cmp.Value != "0.7" || cmp.Rhs != nil {
			t.Fatalf("got %#v", cmp)
		}
	})

	t.Run("SpecRhs", func(t *testing.T) {
		expr, err := parseCompareKo("`기본가` 가 기준 `할인` 의 `최소 가격` 이상이면", 1, "t.md")
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		cmp := expr.(ast.Compare)
		sr, ok := cmp.Rhs.(ast.SpecRef)
		if !ok || len(sr.Steps) != 2 || sr.Steps[0] != "할인" || cmp.Value != "" {
			t.Fatalf("Rhs = %#v", cmp.Rhs)
		}
	})

	t.Run("PathRhsEqualsObjectParticle", func(t *testing.T) {
		expr, err := parseCompareKo("`등급` 이 `작성자` 의 `등급` 과 같으면", 1, "t.md")
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		cmp := expr.(ast.Compare)
		p, ok := cmp.Rhs.(ast.Path)
		if !ok || len(p.Steps) != 2 || cmp.Field != "등급" || cmp.Op != "equals" || cmp.Value != "" {
			t.Fatalf("got %#v", cmp)
		}
	})

	t.Run("CallLhsLiteralRhs", func(t *testing.T) {
		expr, err := parseCompareKo("`nlp`.`score` 의 결과 가 0.7 보다 크면", 1, "t.md")
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		cmp := expr.(ast.Compare)
		c, ok := cmp.Lhs.(ast.Call)
		if !ok || c.Fn.Alias != "nlp" || c.Fn.Name != "score" || cmp.Op != "is greater than" || cmp.Value != "0.7" {
			t.Fatalf("got %#v", cmp)
		}
	})

	t.Run("MapKeyEqualsTrue", func(t *testing.T) {
		// `속성` 의 `내부` 가 true 와 같으면 — map-key path Lhs with a literal true value.
		expr, err := parseCompareKo("`속성` 의 `내부` 가 true 와 같으면", 1, "t.md")
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		cmp := expr.(ast.Compare)
		p, ok := cmp.Lhs.(ast.Path)
		if !ok || len(p.Steps) != 2 || p.Steps[1] != "내부" || cmp.Op != "equals" || cmp.Value != "true" || cmp.Rhs != nil {
			t.Fatalf("got %#v", cmp)
		}
	})

	t.Run("BacktickIsInRhsIsOperand", func(t *testing.T) {
		// Breaking Change (2): a backtick right side is now an operand, not a value.
		expr, err := parseCompareKo("`f` 가 `허용목록` 중 하나면", 1, "t.md")
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		cmp := expr.(ast.Compare)
		if _, ok := cmp.Rhs.(ast.Path); !ok || cmp.Value != "" {
			t.Fatalf("got %#v", cmp)
		}
	})

	t.Run("UnmarkedDotChainLhsErrors", func(t *testing.T) {
		if _, err := parseCompareKo("`a`.`b` 가 100 이상이면", 1, "t.md"); err == nil {
			t.Fatal("expected error: an unmarked dot-chain left side is not a valid operand")
		}
	})
}
