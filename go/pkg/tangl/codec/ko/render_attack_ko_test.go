//ff:func feature=tangl type=parser control=sequence
//ff:what TestRenderAttackKo — renderAttackKo emits `` `Y` 면 `X` 무효다 `` (attacker-first, 받침-selected 면)
package ko

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// TestRenderAttackKo checks the defeat edge form: the attacker leads with a
// 받침-selected 면/이면, then the target and 무효다.
func TestRenderAttackKo(t *testing.T) {
	got := renderAttackKo(ast.Attack{Target: "주문 접수", Attacker: "컵 없음"})
	if got != "`컵 없음` 이면 `주문 접수` 무효다" {
		t.Errorf("renderAttackKo = %q", got)
	}
}
