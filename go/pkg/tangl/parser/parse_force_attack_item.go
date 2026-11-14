//ff:func feature=tangl type=parser control=sequence
//ff:what parseForceAttackItem — parse "`attacker` attacks `target` with force F" (explicit attack force)
package parser

import (
	"strings"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// parseForceAttackItem parses the explicit-force defeat edge
// NAME "attacks" NAME "with" "force" DECIMAL. The force is a weight in [0, 1]
// with at most four decimal places; it is always recorded (Force non-nil, GEUL
// F_FORCE on), even when it is 1. It returns ok=false (no error) when the item
// does not start with a backtick name followed by "attacks".
func parseForceAttackItem(it item, path string) (ast.Attack, bool, error) {
	attacker, rest, ok := takeBacktick(it.Text)
	if !ok {
		return ast.Attack{}, false, nil
	}
	rest, ok = takeKeyword(rest, "attacks")
	if !ok {
		return ast.Attack{}, false, nil
	}
	target, rest, ok := takeBacktick(rest)
	if !ok {
		return ast.Attack{}, true, errAt(path, it.Line, "expected backtick-quoted target after 'attacks'")
	}
	rest, ok = takeKeyword(rest, "with force")
	if !ok {
		return ast.Attack{}, true, errAt(path, it.Line, "expected 'with force' in attacks edge")
	}
	force, rest, err := parseUnitDecimal(rest, "force", path, it.Line)
	if err != nil {
		return ast.Attack{}, true, err
	}
	if strings.TrimSpace(rest) != "" {
		return ast.Attack{}, true, errAt(path, it.Line, "unexpected trailing text: %q", rest)
	}
	return ast.Attack{Target: target, Attacker: attacker, Force: &force, Line: it.Line}, true, nil
}
