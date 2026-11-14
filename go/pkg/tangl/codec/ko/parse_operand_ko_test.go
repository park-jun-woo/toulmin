//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseOperandKo — spec ref / data path / value call discrimination and malformed-operand errors
package ko

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestParseOperandKo(t *testing.T) {
	t.Run("SingleStepPath", func(t *testing.T) {
		op, _, ok, err := parseOperandKo("`기본가`", 1, "t.md")
		if err != nil || !ok {
			t.Fatalf("ok=%v err=%v", ok, err)
		}
		p, isPath := op.(ast.Path)
		if !isPath || len(p.Steps) != 1 || p.Steps[0] != "기본가" {
			t.Fatalf("got %#v", op)
		}
	})

	t.Run("MultiStepPath", func(t *testing.T) {
		op, _, ok, err := parseOperandKo("`작성자` 의 `신뢰점수` 가", 1, "t.md")
		if err != nil || !ok {
			t.Fatalf("ok=%v err=%v", ok, err)
		}
		p, isPath := op.(ast.Path)
		if !isPath || len(p.Steps) != 2 || p.Steps[1] != "신뢰점수" {
			t.Fatalf("got %#v", op)
		}
	})

	t.Run("SpecRef", func(t *testing.T) {
		op, rest, ok, err := parseOperandKo("기준 `할인` 의 `최소 가격`", 1, "t.md")
		if err != nil || !ok {
			t.Fatalf("ok=%v err=%v", ok, err)
		}
		sr, isSpec := op.(ast.SpecRef)
		if !isSpec || len(sr.Steps) != 2 || sr.Steps[0] != "할인" || sr.Steps[1] != "최소 가격" {
			t.Fatalf("got %#v rest=%q", op, rest)
		}
	})

	t.Run("AliasValueCall", func(t *testing.T) {
		op, _, ok, err := parseOperandKo("`nlp`.`score` 의 결과 가", 1, "t.md")
		if err != nil || !ok {
			t.Fatalf("ok=%v err=%v", ok, err)
		}
		c, isCall := op.(ast.Call)
		if !isCall || c.Fn.Alias != "nlp" || c.Fn.Name != "score" {
			t.Fatalf("got %#v", op)
		}
	})

	t.Run("LocalValueCall", func(t *testing.T) {
		op, _, ok, err := parseOperandKo("`점수화` 의 결과 가", 1, "t.md")
		if err != nil || !ok {
			t.Fatalf("ok=%v err=%v", ok, err)
		}
		c, isCall := op.(ast.Call)
		if !isCall || c.Fn.Alias != "" || c.Fn.Name != "점수화" {
			t.Fatalf("got %#v", op)
		}
	})

	t.Run("NotAnOperand", func(t *testing.T) {
		_, _, ok, err := parseOperandKo("650", 1, "t.md")
		if ok || err != nil {
			t.Fatalf("want ok=false err=nil, got ok=%v err=%v", ok, err)
		}
	})

	t.Run("SpecWithoutBacktick", func(t *testing.T) {
		if _, _, _, err := parseOperandKo("기준 이상이면", 1, "t.md"); err == nil {
			t.Fatal("expected error: 기준 without a backtick-quoted spec term")
		}
	})

	t.Run("UnmarkedDotChain", func(t *testing.T) {
		if _, _, _, err := parseOperandKo("`nlp`.`score` 가 0.7", 1, "t.md"); err == nil {
			t.Fatal("expected error: dot-chain FUNC_REF without a 의 결과 marker")
		}
	})

	t.Run("MalformedFuncRef", func(t *testing.T) {
		if _, _, _, err := parseOperandKo("`nlp`.score 의 결과", 1, "t.md"); err == nil {
			t.Fatal("expected error: dot-chain without a second backtick term")
		}
	})

	t.Run("MultiStepPathWithResult", func(t *testing.T) {
		if _, _, _, err := parseOperandKo("`a` 의 `b` 의 결과", 1, "t.md"); err == nil {
			t.Fatal("expected error: a multi-step path cannot take 의 결과")
		}
	})
}
