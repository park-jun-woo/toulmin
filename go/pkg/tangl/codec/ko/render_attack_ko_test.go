//ff:func feature=tangl type=parser control=sequence
//ff:what TestRenderAttackKo — renderAttackKo emits the 무효다 (omitted force) and 세기로 공격한다 (explicit force) forms
package ko

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// TestRenderAttackKo checks both defeat edge forms: the attacker leads with a
// 받침-selected 면/이면, then the target and 무효다 when Force is nil; with a Force the
// target takes a 받침-selected 을/를 and the shortest decimal is followed by
// "의 세기로 공격한다" (1 stays explicit).
func TestRenderAttackKo(t *testing.T) {
	got := renderAttackKo(ast.Attack{Target: "주문 접수", Attacker: "컵 없음"})
	if got != "`컵 없음` 이면 `주문 접수` 무효다" {
		t.Errorf("renderAttackKo = %q", got)
	}
	half, one := 0.5, 1.0
	got = renderAttackKo(ast.Attack{Target: "주문 접수", Attacker: "컵 없음", Force: &half})
	if got != "`컵 없음` 이면 `주문 접수` 를 0.5의 세기로 공격한다" {
		t.Errorf("renderAttackKo(0.5) = %q", got)
	}
	got = renderAttackKo(ast.Attack{Target: "인증", Attacker: "차단", Force: &one})
	if got != "`차단` 이면 `인증` 을 1의 세기로 공격한다" {
		t.Errorf("renderAttackKo(1) = %q", got)
	}
}
