//ff:func feature=tangl type=parser control=sequence
//ff:what parseArgValue — convert a spec call argument string to its typed representation
package parser

import (
	"strconv"
	"strings"
)

// parseArgValue converts a spec call argument string to a typed value.
func parseArgValue(s string) any {
	if (strings.HasPrefix(s, "\"") && strings.HasSuffix(s, "\"")) ||
		(strings.HasPrefix(s, "'") && strings.HasSuffix(s, "'")) {
		return s[1 : len(s)-1]
	}
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f
	}
	if i, err := strconv.Atoi(s); err == nil {
		return i
	}
	if s == "true" {
		return true
	}
	if s == "false" {
		return false
	}
	if s == "nil" {
		return nil
	}
	return s
}
