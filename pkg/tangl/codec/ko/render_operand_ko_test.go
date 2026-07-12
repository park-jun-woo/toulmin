//ff:func feature=tangl type=parser control=sequence
//ff:what TestRenderOperandKo — Path/SpecRef/Call surfaces with their 조사 base and the unsupported-operand error
package ko

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestRenderOperandKo(t *testing.T) {
	t.Run("SingleStepPath", func(t *testing.T) {
		text, base, err := renderOperandKo(ast.Path{Steps: []string{"기본가"}})
		if err != nil || text != "`기본가`" || base != "기본가" {
			t.Fatalf("got text=%q base=%q err=%v", text, base, err)
		}
	})

	t.Run("MultiStepPath", func(t *testing.T) {
		text, base, err := renderOperandKo(ast.Path{Steps: []string{"작성자", "활성"}})
		if err != nil || text != "`작성자` 의 `활성`" || base != "활성" {
			t.Fatalf("got text=%q base=%q err=%v", text, base, err)
		}
	})

	t.Run("SpecRef", func(t *testing.T) {
		text, base, err := renderOperandKo(ast.SpecRef{Steps: []string{"할인", "최소 가격"}})
		if err != nil || text != "기준 `할인` 의 `최소 가격`" || base != "최소 가격" {
			t.Fatalf("got text=%q base=%q err=%v", text, base, err)
		}
	})

	t.Run("Call", func(t *testing.T) {
		text, base, err := renderOperandKo(ast.Call{Fn: ast.Ref{Alias: "nlp", Name: "score"}})
		if err != nil || text != "`nlp`.`score` 의 결과" || base != "결과" {
			t.Fatalf("got text=%q base=%q err=%v", text, base, err)
		}
	})

	t.Run("UnsupportedOperand", func(t *testing.T) {
		if _, _, err := renderOperandKo(&ast.Path{Steps: []string{"x"}}); err == nil {
			t.Fatal("expected error for an unsupported operand type")
		}
	})
}
