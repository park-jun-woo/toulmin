//ff:func feature=tangl type=parser control=sequence
//ff:what parseOperandKo — parse a value operand: 기준 spec ref, 의-chain path, or 의 결과 host call
package ko

import (
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// parseOperandKo parses one value operand and returns it with the remaining text.
// ok reports whether an operand introducer (the 기준 keyword or a leading backtick)
// was recognized; when ok is false the caller decides how to treat the text (e.g. as
// a literal value). A non-nil error means an introducer was recognized but the
// operand was malformed. It is the KO mirror of pkg/tangl/parser.parseOperand: 기준
// opens a SpecRef, a dot-chained backtick (`a`.`b`) is a FUNC_REF that must carry the
// 의 결과 value-call marker, and a plain 의-chain is a data Path (a single-step path
// followed by 의 결과 is a bare local host call).
func parseOperandKo(s string, line int, path string) (ast.Operand, string, bool, error) {
	if rest, ok := takeKeywordKo(s, "기준"); ok {
		steps, r2, ok2 := parsePathKo(rest)
		if !ok2 {
			return nil, s, true, errAtKo(path, line, "expected a backtick-quoted spec term after '기준'")
		}
		return ast.SpecRef{Steps: steps}, r2, true, nil
	}
	_, after1, ok1 := takeBacktick(s)
	if !ok1 {
		return nil, s, false, nil
	}
	if len(after1) > 0 && after1[0] == '.' {
		ref, rest, ok2 := parseRefKo(s)
		if !ok2 {
			return nil, s, true, errAtKo(path, line, "malformed FUNC_REF in %q", s)
		}
		if r3, ok3 := takeKeywordKo(skipParticlesKo(rest), "결과"); ok3 {
			return ast.Call{Fn: ref}, r3, true, nil
		}
		return nil, s, true, errAtKo(path, line, "a dot-chain FUNC_REF needs a '의 결과' or '참이면' marker in %q", s)
	}
	steps, rest, ok3 := parsePathKo(s)
	if !ok3 {
		return nil, s, false, nil
	}
	if r3, ok4 := takeKeywordKo(skipParticlesKo(rest), "결과"); ok4 {
		if len(steps) != 1 {
			return nil, s, true, errAtKo(path, line, "a '의 결과' host call requires a `별칭`.`이름` FUNC_REF, not a %d-step path", len(steps))
		}
		return ast.Call{Fn: ast.Ref{Name: steps[0]}}, r3, true, nil
	}
	return ast.Path{Steps: steps}, rest, true, nil
}
