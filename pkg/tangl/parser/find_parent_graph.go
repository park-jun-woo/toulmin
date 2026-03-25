//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what findParentGraph — search backwards through lines to find the enclosing graph name
package parser

import "strings"

// findParentGraph searches backwards from current line to find the parent graph declaration.
func findParentGraph(lines []string, current int, currentIndent int) string {
	for j := current - 1; j >= 0; j-- {
		raw := lines[j]
		trimmed := strings.TrimSpace(raw)
		if trimmed == "" {
			continue
		}
		ind := countIndent(raw)
		if ind >= currentIndent {
			continue
		}
		stripped := stripListPrefix(trimmed)
		if !strings.Contains(stripped, "is a graph ") {
			return ""
		}
		parts := strings.SplitN(stripped, " is a graph ", 2)
		if len(parts) == 2 {
			return parts[0]
		}
		return ""
	}
	return ""
}
