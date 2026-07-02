//ff:func feature=tangl type=parser control=sequence
//ff:what parseRequireItem — parse a "`field` is required [as Type]" declaration
package parser

import (
	"strings"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// parseRequireItem parses "`<field>` is required [as <Type>]". It returns
// ok=false (no error) when the item does not start this way, so the caller
// can try the next statement form (e.g. a node registration).
func parseRequireItem(it item, path string) (ast.Require, bool, error) {
	field, rest, ok := takeBacktick(it.Text)
	if !ok {
		return ast.Require{}, false, nil
	}
	rest, ok = takeKeyword(rest, "is")
	if !ok {
		return ast.Require{}, false, nil
	}
	rest, ok = takeKeyword(rest, "required")
	if !ok {
		return ast.Require{}, false, nil
	}
	req := ast.Require{Field: field, Line: it.Line}
	rest = strings.TrimSpace(rest)
	if rest == "" {
		return req, true, nil
	}
	rest, ok = takeKeyword(rest, "as")
	if !ok {
		return ast.Require{}, true, errAt(path, it.Line, "unexpected trailing text: %q", rest)
	}
	typ, rest, ok := parseTypeRef(rest)
	if !ok {
		return ast.Require{}, true, errAt(path, it.Line, "expected type after 'as'")
	}
	if strings.TrimSpace(rest) != "" {
		return ast.Require{}, true, errAt(path, it.Line, "unexpected trailing text: %q", rest)
	}
	req.Type = typ
	return req, true, nil
}
