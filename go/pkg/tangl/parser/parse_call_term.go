//ff:func feature=tangl type=parser control=sequence
//ff:what tryParseCallTerm — parse a "FUNC_REF holds" predicate term, falling through otherwise
package parser

import (
	"strings"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// tryParseCallTerm parses a host predicate term of the form "FUNC_REF holds".
// It returns consumed=false (no error) when the text is not a call term, so
// the caller can fall through to parseCompare (backtick starts are shared with
// comparisons). A non-nil error is a malformed call term (trailing text).
func tryParseCallTerm(text string, line int, path string) (ast.Expr, bool, error) {
	ref, rest, ok := parseRef(text)
	if !ok {
		return nil, false, nil
	}
	rest, ok = takeKeyword(rest, "holds")
	if !ok {
		return nil, false, nil
	}
	if strings.TrimSpace(rest) != "" {
		return nil, true, errAt(path, line, "unexpected trailing text after 'holds': %q", rest)
	}
	return ast.Call{Fn: ref}, true, nil
}
