//ff:func feature=tangl type=parser control=sequence
//ff:what parseFieldItem — parse a StructDef field ("has `f` as Type")
package parser

import (
	"strings"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// parseFieldItem parses "has `<field>` as <Type>".
func parseFieldItem(it item, path string) (ast.Field, error) {
	rest, ok := takeKeyword(it.Text, "has")
	if !ok {
		return ast.Field{}, errAt(path, it.Line, "expected 'has `field` as Type', got %q", it.Text)
	}
	name, rest, ok := takeBacktick(rest)
	if !ok {
		return ast.Field{}, errAt(path, it.Line, "expected backtick-quoted field name after 'has'")
	}
	rest, ok = takeKeyword(rest, "as")
	if !ok {
		return ast.Field{}, errAt(path, it.Line, "expected 'as Type' after field name")
	}
	typ, rest, ok := parseTypeRef(rest)
	if !ok {
		return ast.Field{}, errAt(path, it.Line, "expected type after 'as'")
	}
	if strings.TrimSpace(rest) != "" {
		return ast.Field{}, errAt(path, it.Line, "unexpected trailing text: %q", rest)
	}
	return ast.Field{Name: name, Type: typ, Line: it.Line}, nil
}
