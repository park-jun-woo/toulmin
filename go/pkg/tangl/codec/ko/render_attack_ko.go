//ff:func feature=tangl type=parser control=sequence
//ff:what renderAttackKo — render a defeat edge (`` `Y` 면 `X` 무효다 `` or `` `Y` 면 `X` 를 F의 세기로 공격한다 ``)
package ko

import (
	"strconv"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// renderAttackKo renders a defeat edge back to "`<attacker>` 면 `<target>` 무효다"
// when Force is nil (omitted, force 1), or "`<attacker>` 면 `<target>` 를 <F>의 세기로
// 공격한다" when Force is set (explicit, even for 1). Unlike the English forms, the
// KO 정규형 puts the attacker first and the target second (the inverse of
// parseAttackKo/parseForceAttackKo). The 면/이면 and 을/를 particles are
// 받침-selected; F is the shortest decimal, the authored spelling surviving only
// in mode-1 (SOURCE) decoding.
func renderAttackKo(a ast.Attack) string {
	head := "`" + a.Attacker + "` " + josaFor(a.Attacker, "이면", "면") + " `" + a.Target + "` "
	if a.Force == nil {
		return head + "무효다"
	}
	return head + josaFor(a.Target, "을", "를") + " " + strconv.FormatFloat(*a.Force, 'g', -1, 64) + "의 세기로 공격한다"
}
