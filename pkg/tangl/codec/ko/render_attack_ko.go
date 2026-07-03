//ff:func feature=tangl type=parser control=sequence
//ff:what renderAttackKo — render a defeat edge (`` `Y` 면 `X` 무효다 ``, attacker-first)
package ko

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// renderAttackKo renders a defeat edge back to "`<attacker>` 면 `<target>` 무효다". Unlike
// the English "don't `<target>` when `<attacker>`", the KO 정규형 puts the attacker
// first and the target second (the inverse of parseAttackKo). The 면(면/이면) particle
// after the attacker is 받침-selected.
func renderAttackKo(a ast.Attack) string {
	return "`" + a.Attacker + "` " + josaFor(a.Attacker, "이면", "면") + " `" + a.Target + "` 무효다"
}
