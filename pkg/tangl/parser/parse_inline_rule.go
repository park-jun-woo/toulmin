//ff:func feature=tangl type=parser control=sequence
//ff:what parseInlineRule — parse an inline rule header and body into InlineRule AST node
package parser

import (
	"fmt"
	"strings"
)

// parseInlineRule parses header: rule "name" is, body: return that <expr>
func parseInlineRule(header string, body string, lineNum int) (InlineRule, error) {
	nameStart := strings.Index(header, "\"")
	nameEnd := strings.LastIndex(header, "\"")
	if nameStart < 0 || nameEnd <= nameStart {
		return InlineRule{}, fmt.Errorf("invalid inline rule header: %s", header)
	}
	name := header[nameStart+1 : nameEnd]

	body = strings.TrimSpace(body)
	if !strings.HasPrefix(body, "return that ") {
		return InlineRule{}, fmt.Errorf("inline rule body must start with 'return that': %s", body)
	}
	exprText := strings.TrimPrefix(body, "return that ")

	expr, err := parseExpr(exprText)
	if err != nil {
		return InlineRule{}, fmt.Errorf("invalid expression in inline rule %q: %w", name, err)
	}

	return InlineRule{Name: name, Expr: expr, Line: lineNum}, nil
}
