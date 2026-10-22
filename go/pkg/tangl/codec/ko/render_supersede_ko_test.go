//ff:func feature=tangl type=parser control=sequence
//ff:what TestRenderSupersedeKo — renderSupersedeKo emits `` `X` 면 `Y` 배제한다 `` (superseder-first, 받침-selected 면)
package ko

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// TestRenderSupersedeKo checks the supersede edge form: the superseder leads with a
// 받침-selected 면/이면 (no order reversal), then the superseded target and 배제한다.
func TestRenderSupersedeKo(t *testing.T) {
	t.Run("WithBatchim", func(t *testing.T) {
		// 없음 ends with a 받침 (ㅁ) → 이면
		got := renderSupersedeKo(ast.Supersede{Superseder: "컵 없음", Superseded: "주문 접수"})
		if got != "`컵 없음` 이면 `주문 접수` 배제한다" {
			t.Errorf("renderSupersedeKo = %q", got)
		}
	})

	t.Run("WithoutBatchim", func(t *testing.T) {
		// 접수 ends without a 받침 (수) → 면
		got := renderSupersedeKo(ast.Supersede{Superseder: "주문 접수", Superseded: "컵 없음"})
		if got != "`주문 접수` 면 `컵 없음` 배제한다" {
			t.Errorf("renderSupersedeKo = %q", got)
		}
	})
}
