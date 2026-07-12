//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseCompareRhsKo — literal Value fallback, operand Rhs, trailing-text and malformed-operand errors
package ko

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestParseCompareRhsKo(t *testing.T) {
	seed := ast.Compare{Field: "점수", Op: "is at least"}

	t.Run("LiteralValue", func(t *testing.T) {
		cmp, err := parseCompareRhsKo(seed, " 가 650 와", 1, "t.md")
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if cmp.Value != "650" || cmp.Rhs != nil {
			t.Fatalf("got %#v", cmp)
		}
	})

	t.Run("OperandRhs", func(t *testing.T) {
		cmp, err := parseCompareRhsKo(seed, " 이 `기본가`", 1, "t.md")
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		p, isPath := cmp.Rhs.(ast.Path)
		if cmp.Value != "" || !isPath || len(p.Steps) != 1 || p.Steps[0] != "기본가" {
			t.Fatalf("got %#v", cmp)
		}
	})

	t.Run("TrailingTextAfterOperand", func(t *testing.T) {
		if _, err := parseCompareRhsKo(seed, "`기본가` 잔여", 1, "t.md"); err == nil {
			t.Fatal("expected error: residual bytes past the operand")
		}
	})

	t.Run("MalformedOperand", func(t *testing.T) {
		if _, err := parseCompareRhsKo(seed, "기준 뒤", 1, "t.md"); err == nil {
			t.Fatal("expected error: 기준 without a backtick-quoted spec term")
		}
	})
}
