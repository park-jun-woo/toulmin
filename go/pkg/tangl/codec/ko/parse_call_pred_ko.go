//ff:func feature=tangl type=parser control=sequence
//ff:what parseCallPredKo — parse a `FUNC_REF 가 참이면` host predicate term, falling through otherwise
package ko

import (
	"strings"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// parseCallPredKo parses a host predicate term of the form `FUNC_REF 가 참이면`. It
// returns consumed=false (no error) when the text does not end in 참이면, so the caller
// can fall through to a comparison or bare name. Once the 참이면 suffix is committed the
// leading subject 조사 is stripped and the FUNC_REF read; residual non-particle text is
// an error. It is the KO mirror of pkg/tangl/parser.tryParseCallTerm (whose marker is
// the trailing 'holds'); in earlier KO documents a term ending in 참이면 was always an
// "unrecognized condition" error, so this branch is purely additive.
func parseCallPredKo(text string, line int, path string) (ast.Expr, bool, error) {
	before, ok := takeSuffixKo(text, "참이면")
	if !ok {
		return nil, false, nil
	}
	ref, rest, ok := parseRefKo(trimTrailingParticleKo(before))
	if !ok {
		return nil, true, errAtKo(path, line, "expected a FUNC_REF before '참이면' in %q", text)
	}
	if strings.TrimSpace(skipParticlesKo(rest)) != "" {
		return nil, true, errAtKo(path, line, "unexpected trailing text before '참이면': %q", rest)
	}
	return ast.Call{Fn: ref}, true, nil
}
