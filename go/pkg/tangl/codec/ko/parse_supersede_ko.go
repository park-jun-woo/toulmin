//ff:func feature=tangl type=parser control=sequence
//ff:what parseSupersedeKo — parse '`X` 면 `Y` 배제한다' into an ast.Supersede (no order swap)
package ko

import (
	"strings"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/codec/scan"
)

// parseSupersedeKo parses the supersede edge "`X` 면 `Y` 배제한다". Unlike parseAttackKo
// (which reverses order so the attacker leads), Supersede does NOT reverse: the FIRST
// backtick is the Superseder (X) and the SECOND is the Superseded target (Y), exactly
// as in the English "`X` supersedes `Y`". It returns ok=false (no error) when the item
// does not end in 배제한다.
func parseSupersedeKo(it scan.Item, path string) (ast.Supersede, bool, error) {
	body, ok := takeSuffixKo(it.Text, "배제한다")
	if !ok {
		return ast.Supersede{}, false, nil
	}
	superseder, rest, ok := takeBacktick(body)
	if !ok {
		return ast.Supersede{}, true, errAtKo(path, it.No, "expected backtick-quoted superseder before '배제한다'")
	}
	rest = skipParticlesKo(rest)
	superseded, rest, ok := takeBacktick(rest)
	if !ok {
		return ast.Supersede{}, true, errAtKo(path, it.No, "expected backtick-quoted superseded target in '배제한다' edge")
	}
	if strings.TrimSpace(skipParticlesKo(rest)) != "" {
		return ast.Supersede{}, true, errAtKo(path, it.No, "unexpected trailing text: %q", rest)
	}
	return ast.Supersede{Superseder: superseder, Superseded: superseded, Line: it.No}, true, nil
}
