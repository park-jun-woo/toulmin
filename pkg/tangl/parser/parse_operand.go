//ff:func feature=tangl type=parser control=sequence
//ff:what parseOperand — parse a value operand: path, spec reference, or host call
package parser

import (
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// parseOperand parses one value operand and returns it with the remaining
// text. ok reports whether an operand introducer (a backtick, the "spec"
// keyword, or the "result" keyword) was recognized; when ok is false the
// caller decides how to treat the text (e.g. as a literal value). A non-nil
// error means an introducer was recognized but the operand was malformed.
func parseOperand(s string, line int, path string) (ast.Operand, string, bool, error) {
	if rest, ok := takeKeyword(s, "spec"); ok {
		p, r2, ok2 := parsePath(rest)
		if !ok2 {
			return nil, s, true, errAt(path, line, "expected backtick-quoted spec term after 'spec'")
		}
		return ast.SpecRef{Steps: p.Steps}, r2, true, nil
	}
	if rest, ok := takeKeyword(s, "result"); ok {
		rest, ok = takeKeyword(rest, "of")
		if !ok {
			return nil, s, true, errAt(path, line, "expected 'of' after 'result'")
		}
		ref, r2, ok2 := parseRef(rest)
		if !ok2 {
			return nil, s, true, errAt(path, line, "expected backtick-quoted function reference after 'result of'")
		}
		return ast.Call{Fn: ref}, r2, true, nil
	}
	if p, r2, ok := parsePath(s); ok {
		return p, r2, true, nil
	}
	return nil, s, false, nil
}
