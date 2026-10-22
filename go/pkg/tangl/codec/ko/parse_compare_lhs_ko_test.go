//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseCompareLhsKo — single-step Field demotion, multi-step/spec Lhs, fall-through, and errors
package ko

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestParseCompareLhsKo(t *testing.T) {
	t.Run("SingleStepToField", func(t *testing.T) {
		cmp, rest, ok, err := parseCompareLhsKo("`점수` 가", 1, "t.md")
		if err != nil || !ok {
			t.Fatalf("ok=%v err=%v", ok, err)
		}
		if cmp.Field != "점수" || cmp.Lhs != nil || rest != " 가" {
			t.Fatalf("got %#v rest=%q", cmp, rest)
		}
	})

	t.Run("MultiStepToLhs", func(t *testing.T) {
		cmp, _, ok, err := parseCompareLhsKo("`c` 의 `최소가격` 이", 1, "t.md")
		if err != nil || !ok {
			t.Fatalf("ok=%v err=%v", ok, err)
		}
		p, isPath := cmp.Lhs.(ast.Path)
		if cmp.Field != "" || !isPath || len(p.Steps) != 2 || p.Steps[1] != "최소가격" {
			t.Fatalf("got %#v", cmp)
		}
	})

	t.Run("SpecRefToLhs", func(t *testing.T) {
		cmp, _, ok, err := parseCompareLhsKo("기준 `할인` 의 `최소 가격` 이", 1, "t.md")
		if err != nil || !ok {
			t.Fatalf("ok=%v err=%v", ok, err)
		}
		sr, isSpec := cmp.Lhs.(ast.SpecRef)
		if !isSpec || len(sr.Steps) != 2 || sr.Steps[0] != "할인" {
			t.Fatalf("got %#v", cmp)
		}
	})

	t.Run("NoOperandFallThrough", func(t *testing.T) {
		_, rest, ok, err := parseCompareLhsKo("650", 1, "t.md")
		if ok || err != nil || rest != "650" {
			t.Fatalf("want ok=false err=nil rest=%q, got ok=%v err=%v rest=%q", "650", ok, err, rest)
		}
	})

	t.Run("MalformedOperand", func(t *testing.T) {
		if _, _, _, err := parseCompareLhsKo("기준 없음", 1, "t.md"); err == nil {
			t.Fatal("expected error: 기준 without a backtick-quoted spec term")
		}
	})
}
