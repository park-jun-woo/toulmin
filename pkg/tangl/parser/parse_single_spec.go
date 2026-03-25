//ff:func feature=tangl type=parser control=sequence
//ff:what parseSingleSpec — parse a single spec call like FuncName("arg1", "arg2") into SpecCall
package parser

import (
	"fmt"
	"strings"
)

// parseSingleSpec parses a single spec call string into a SpecCall.
func parseSingleSpec(text string) (SpecCall, error) {
	parenIdx := strings.Index(text, "(")
	if parenIdx < 0 {
		return SpecCall{}, fmt.Errorf("invalid spec call (no parenthesis): %s", text)
	}
	name := strings.TrimSpace(text[:parenIdx])
	argsStr := text[parenIdx+1:]
	closeIdx := strings.LastIndex(argsStr, ")")
	if closeIdx < 0 {
		return SpecCall{}, fmt.Errorf("invalid spec call (no closing parenthesis): %s", text)
	}
	argsStr = argsStr[:closeIdx]

	var args []any
	if strings.TrimSpace(argsStr) != "" {
		rawArgs := strings.Split(argsStr, ",")
		for _, a := range rawArgs {
			a = strings.TrimSpace(a)
			args = append(args, parseArgValue(a))
		}
	}

	return SpecCall{Name: name, Args: args}, nil
}
