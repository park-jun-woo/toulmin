//ff:func feature=tangl type=parser control=sequence
//ff:what parseRequireKo — parse '`f` 는 필수다[, 타입 T]' into an ast.Require
package ko

import (
	"strings"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/codec/scan"
)

// parseRequireKo parses "`field` 는 필수다[, 타입 <Type>]". 필수다 sits mid-statement
// (a ', 타입 <T>' tail may follow), so it is located by contains-keyword dispatch.
// It returns ok=false (no error) when the item is not a require declaration, so the
// caller can try the next case-body form.
func parseRequireKo(it scan.Item, path string) (ast.Require, bool, error) {
	field, rest, ok := takeBacktick(it.Text)
	if !ok {
		return ast.Require{}, false, nil
	}
	rest = skipParticlesKo(rest)
	idx := strings.Index(rest, "필수다")
	if idx < 0 {
		return ast.Require{}, false, nil
	}
	if strings.TrimSpace(rest[:idx]) != "" {
		return ast.Require{}, false, nil
	}
	req := ast.Require{Field: field, Line: it.No}
	tail := strings.TrimSpace(rest[idx+len("필수다"):])
	if tail == "" {
		return req, true, nil
	}
	t := skipParticlesKo(tail)
	t, ok = takeKeywordKo(t, "타입")
	if !ok {
		return ast.Require{}, true, errAtKo(path, it.No, "unexpected trailing text after '필수다': %q", tail)
	}
	typ, rest2, ok := parseTypeRefKo(strings.TrimSpace(t))
	if !ok {
		return ast.Require{}, true, errAtKo(path, it.No, "expected type after '타입'")
	}
	if strings.TrimSpace(rest2) != "" {
		return ast.Require{}, true, errAtKo(path, it.No, "unexpected trailing text: %q", rest2)
	}
	req.Type = typ
	return req, true, nil
}
