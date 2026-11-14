//ff:func feature=tangl type=parser control=sequence
//ff:what renderAttack — render a defeat edge (omitted force: don't `target` when `attacker`; explicit: `attacker` attacks `target` with force F)
package en

import (
	"strconv"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// renderAttack renders a defeat edge back to its canonical form: "don't
// `<target>` when `<attacker>`" when Force is nil (omitted, force 1), or
// "`<attacker>` attacks `<target>` with force <F>" when Force is set (explicit,
// even for 1). F is the shortest decimal ("0.5", "1"); the authored spelling
// ("0.50", "1.0") survives only in mode-1 (SOURCE) decoding.
func renderAttack(a ast.Attack) string {
	if a.Force == nil {
		return "don't `" + a.Target + "` when `" + a.Attacker + "`"
	}
	return "`" + a.Attacker + "` attacks `" + a.Target + "` with force " + strconv.FormatFloat(*a.Force, 'g', -1, 64)
}
