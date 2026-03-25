//ff:func feature=tangl type=parser control=sequence
//ff:what parseSingleExpr — parse a single comparison expression into Expr
package parser

import (
	"fmt"
	"strings"
)

// parseSingleExpr parses a single comparison expression (field operator value).
func parseSingleExpr(text string) (Expr, error) {
	e := Expr{}
	words := strings.Fields(text)
	if len(words) < 2 {
		return Expr{}, fmt.Errorf("expression too short: %s", text)
	}

	e.Field = words[0]
	rest := words[1:]

	if len(rest) >= 2 && rest[0] == "of" && rest[1] == "spec" {
		e.OfSpec = true
		rest = rest[2:]
	}

	op, consumed, err := matchOperator(rest)
	if err != nil {
		return Expr{}, fmt.Errorf("no operator in expression: %s", text)
	}
	e.Operator = op
	rest = rest[consumed:]

	if op == "is nil" || op == "is not nil" {
		return e, nil
	}

	if len(rest) == 0 {
		return Expr{}, fmt.Errorf("missing value after operator in: %s", text)
	}

	valStr := strings.Join(rest, " ")
	e.Value = parseExprValue(valStr)

	return e, nil
}
