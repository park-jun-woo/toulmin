//ff:func feature=tangl type=parser control=sequence
//ff:what parseAttackKo — parse '`Y` 면 `X` 무효다' into an ast.Attack (attacker-first order)
package ko

import (
	"strings"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/pkg/tangl/codec/scan"
)

// parseAttackKo parses the defeat edge "`Y` 면 `X` 무효다". Unlike the English
// "don't `X` when `Y`", the KO 정규형 puts the attacker first and the target second,
// so the FIRST backtick is the Attacker and the SECOND is the Target. It returns
// ok=false (no error) when the item does not end in 무효다.
func parseAttackKo(it scan.Item, path string) (ast.Attack, bool, error) {
	body, ok := takeSuffixKo(it.Text, "무효다")
	if !ok {
		return ast.Attack{}, false, nil
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
