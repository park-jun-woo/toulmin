//ff:func feature=tangl type=parser control=sequence
//ff:what parseAttackKo — parse '`Y` 면 `X` 무효다' (or the explicit-force '… 공격한다') into an ast.Attack
package ko

import (
	"strings"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/codec/scan"
)

// parseAttackKo parses the defeat edge "`Y` 면 `X` 무효다". Unlike the English
// "don't `X` when `Y`", the KO 정규형 puts the attacker first and the target second,
// so the FIRST backtick is the Attacker and the SECOND is the Target. An item that
// does not end in 무효다 is handed to parseForceAttackKo (the explicit-force form
// '`Y` 면 `X` 를 F의 세기로 공격한다'), which returns ok=false when it is not that
// either.
func parseAttackKo(it scan.Item, path string) (ast.Attack, bool, error) {
	body, ok := takeSuffixKo(it.Text, "무효다")
	if !ok {
		return parseForceAttackKo(it, path)
	}
	attacker, rest, ok := takeBacktick(body)
	if !ok {
		return ast.Attack{}, true, errAtKo(path, it.No, "expected backtick-quoted attacker before '무효다'")
	}
	rest = skipParticlesKo(rest)
	target, rest, ok := takeBacktick(rest)
	if !ok {
		return ast.Attack{}, true, errAtKo(path, it.No, "expected backtick-quoted target in '무효다' edge")
	}
	if strings.TrimSpace(skipParticlesKo(rest)) != "" {
		return ast.Attack{}, true, errAtKo(path, it.No, "unexpected trailing text: %q", rest)
	}
	return ast.Attack{Target: target, Attacker: attacker, Line: it.No}, true, nil
}
