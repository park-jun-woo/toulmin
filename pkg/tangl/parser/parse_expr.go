//ff:func feature=tangl type=parser control=sequence
//ff:what parseExpr — parse an inline rule expression into Expr AST node
package parser

import (
	"fmt"
	"strings"
)

// parseExpr parses: field [of spec] operator value [and/or expr]
func parseExpr(text string) (Expr, error) {
	text = strings.TrimSpace(text)
	if text == "" {
		return Expr{}, fmt.Errorf("empty expression")
	}

	andIdx := findLogicalOp(text, " and ")
	if andIdx >= 0 {
		left, err := parseExpr(text[:andIdx])
		if err != nil {
			return Expr{}, err
		}
		right, err := parseExpr(text[andIdx+5:])
		if err != nil {
			return Expr{}, err
		}
		left.And = &right
		return left, nil
	}

	orIdx := findLogicalOp(text, " or ")
	if orIdx >= 0 {
		left, err := parseExpr(text[:orIdx])
		if err != nil {
			return Expr{}, err
		}
		right, err := parseExpr(text[orIdx+4:])
		if err != nil {
			return Expr{}, err
		}
		left.Or = &right
		return left, nil
	}

	return parseSingleExpr(text)
}
