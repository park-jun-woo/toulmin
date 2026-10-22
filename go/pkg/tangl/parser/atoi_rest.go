//ff:func feature=tangl type=parser control=sequence
//ff:what atoiRest — parse a leading integer of known width and return the rest
package parser

import "strconv"

// atoiRest parses the first width bytes of s as an integer and returns the
// remaining text after it.
func atoiRest(s string, width int, path string, line int) (int, string, error) {
	n, err := strconv.Atoi(s[:width])
	if err != nil {
		return 0, "", errAt(path, line, "invalid percent %q: %v", s[:width], err)
	}
	return n, s[width:], nil
}
