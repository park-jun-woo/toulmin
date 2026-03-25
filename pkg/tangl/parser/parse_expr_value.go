//ff:func feature=tangl type=parser control=sequence
//ff:what parseExprValue — convert a string value to its typed representation
package parser

import (
	"strconv"
	"strings"
)

// parseExprValue converts a string to a typed value (int, float, bool, string, or nil).
func parseExprValue(s string) any {
	s = strings.TrimSpace(s)
	if s == "nil" {
		return nil
	}
	if (strings.HasPrefix(s, "\"") && strings.HasSuffix(s, "\"")) ||
		(strings.HasPrefix(s, "'") && strings.HasSuffix(s, "'")) {
		return s[1 : len(s)-1]
	}
	if i, err := strconv.Atoi(s); err == nil {
		return i
	}
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f
	}
	if s == "true" {
		return true
	}
	if s == "false" {
		return false
	}
	return s
}
