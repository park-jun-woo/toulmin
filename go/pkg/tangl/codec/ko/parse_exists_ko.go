//ff:func feature=tangl type=parser control=sequence pattern=recursive
//ff:what parseExistsKo — parse an `<operand> 중 `c` 가 …것이 있으면` exists term (one-line or 다음인 nested)
package ko

import (
	"strings"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/codec/scan"
)

// parseExistsKo tries the exists lookahead: a source operand, the bare 중 connective,
// and a backtick binder in immediate succession. Any miss returns ok=false with no
// error so the caller falls through to a comparison or bare name (the operand's own
// error is swallowed here — a bare dot-chain, for instance, is really a 참이면 call or
// a plain comparison and reports its error there). Once committed, the 것이 있으면 tail
// is required; a body of exactly 다음인 opens the nested child condition list, and any
// other body is one attributive comparison. It is the KO mirror of
// pkg/tangl/parser.parseExistsTerm.
func parseExistsKo(text string, children []scan.Item, line int, path string) (ast.Expr, bool, error) {
	src, rest, ok, err := parseOperandKo(text, line, path)
	if err != nil || !ok {
		return nil, false, nil
	}
	rest, ok = takeKeywordKo(rest, "중")
	if !ok {
		return nil, false, nil
	}
	binder, rest, ok := takeBacktick(rest)
	if !ok {
		return nil, false, nil
	}
	body, ok := takeSuffixKo(skipParticlesKo(rest), "것이 있으면")
	if !ok {
		return nil, true, errAtKo(path, line, "expected '것이 있으면' to close the exists term in %q", text)
	}
	if strings.TrimSpace(body) == "다음인" {
		cond, err := parseCondListKo(children, path)
		if err != nil {
			return nil, true, err
		}
		return ast.Exists{Source: src, Binder: binder, Cond: cond}, true, nil
	}
	cond, err := parseCompareAdnKo(body, line, path)
	if err != nil {
		return nil, true, err
	}
	return ast.Exists{Source: src, Binder: binder, Cond: cond}, true, nil
}
