//ff:func feature=tangl type=parser control=sequence
//ff:what parseConstValueKo — split a ConstDef's literal value from its ', 타입 <Ref>' tail
package ko

import (
	"strings"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// parseConstValueKo turns the value head "<값> 을" and the optional tail
// ", 타입 <Ref>" of a 뜻한다 definition into the raw literal text and an optional
// spec Ref. The value keeps its verbatim literal source (e.g. "650", "65°C"); the
// trailing object particle (을/를) is dropped.
func parseConstValueKo(valuePart, tail string, line int, path string) (string, *ast.Ref, error) {
	value := trimTrailingParticleKo(valuePart)
	if tail == "" {
		return value, nil, nil
	}
	t := skipParticlesKo(tail)
	t, ok := takeKeywordKo(t, "타입")
	if !ok {
		return "", nil, errAtKo(path, line, "expected '타입' in definition tail, got %q", tail)
	}
	ref, rest, ok := parseRefKo(strings.TrimSpace(t))
	if !ok {
		return "", nil, errAtKo(path, line, "expected reference after '타입'")
	}
	if strings.TrimSpace(rest) != "" {
		return "", nil, errAtKo(path, line, "unexpected trailing text: %q", rest)
	}
	return value, &ref, nil
}
