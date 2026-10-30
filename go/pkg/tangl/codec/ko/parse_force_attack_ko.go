//ff:func feature=tangl type=parser control=sequence
//ff:what parseForceAttackKo — parse '`Y` 면 `X` 를 F의 세기로 공격한다' (explicit attack force) into an ast.Attack
package ko

import (
	"strings"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/codec/scan"
)

// parseForceAttackKo parses the explicit-force defeat edge
// "`Y` 면 `X` 를 <F>의 세기로 공격한다": attacker first, target second (as in 무효다),
// then the force F — a weight in [0, 1] with at most four decimal places — its
// '의' (attached or spaced), and the keyword pair 세기로 공격한다. The force is always
// recorded (Force non-nil, GEUL F_FORCE on), even when it is 1. It returns
// ok=false (no error) when the item does not end in 공격한다.
func parseForceAttackKo(it scan.Item, path string) (ast.Attack, bool, error) {
	body, ok := takeSuffixKo(it.Text, "공격한다")
	if !ok {
		return ast.Attack{}, false, nil
	}
	body, ok = takeSuffixKo(body, "세기로")
	if !ok {
		return ast.Attack{}, true, errAtKo(path, it.No, "expected '세기로' before '공격한다'")
	}
	attacker, rest, ok := takeBacktick(body)
	if !ok {
		return ast.Attack{}, true, errAtKo(path, it.No, "expected backtick-quoted attacker before '공격한다'")
	}
	target, rest, ok := takeBacktick(skipParticlesKo(rest))
	if !ok {
		return ast.Attack{}, true, errAtKo(path, it.No, "expected backtick-quoted target in '공격한다' edge")
	}
	force, rest, err := parseUnitDecimalKo(skipParticlesKo(rest), "세기", it.No, path)
	if err != nil {
		return ast.Attack{}, true, err
	}
	if strings.TrimSpace(skipParticlesKo(rest)) != "" {
		return ast.Attack{}, true, errAtKo(path, it.No, "unexpected trailing text before '세기로': %q", rest)
	}
	return ast.Attack{Target: target, Attacker: attacker, Force: &force, Line: it.No}, true, nil
}
