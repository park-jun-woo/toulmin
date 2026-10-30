//ff:func feature=tangl type=parser control=sequence
//ff:what parseAttackItem — parse a "don't/do not `target` when `attacker`" or "`attacker` attacks `target` with force F" edge
package parser

import (
	"strings"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// parseAttackItem parses edge_decl: DONT NAME "when" NAME, where DONT is
// "don't" or "do not" (the omitted-force form, Force nil), or delegates to
// parseForceAttackItem for the explicit form NAME "attacks" NAME "with" "force"
// DECIMAL. It returns ok=false (no error) when the item is neither, so the
// caller can try the exec edge forms next.
func parseAttackItem(it item, path string) (ast.Attack, bool, error) {
	rest, ok := takeKeyword(it.Text, "don't")
	if !ok {
		r2, ok2 := takeKeyword(it.Text, "do")
		if !ok2 {
			return parseForceAttackItem(it, path)
		}
		r3, ok3 := takeKeyword(r2, "not")
		if !ok3 {
			return ast.Attack{}, false, nil
		}
		rest = r3
	}
	target, rest, ok := takeBacktick(rest)
	if !ok {
		return ast.Attack{}, true, errAt(path, it.Line, "expected backtick-quoted target after don't/do not")
	}
	rest, ok = takeKeyword(rest, "when")
	if !ok {
		return ast.Attack{}, true, errAt(path, it.Line, "expected 'when' in don't edge")
	}
	attacker, rest, ok := takeBacktick(rest)
	if !ok {
		return ast.Attack{}, true, errAt(path, it.Line, "expected backtick-quoted attacker after 'when'")
	}
	if strings.TrimSpace(rest) != "" {
		return ast.Attack{}, true, errAt(path, it.Line, "unexpected trailing text: %q", rest)
	}
	return ast.Attack{Target: target, Attacker: attacker, Line: it.Line}, true, nil
}
