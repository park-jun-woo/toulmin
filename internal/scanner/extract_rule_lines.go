//ff:func feature=scanner type=scanner control=iteration dimension=1
//ff:what extractRuleLines — filters //rule: prefixed lines from a comment group
package scanner

import (
	"go/ast"
	"strings"
)

// extractRuleLines returns comment texts with //rule: prefix from a CommentGroup.
func extractRuleLines(doc *ast.CommentGroup) []string {
	var lines []string
	for _, c := range doc.List {
		if strings.HasPrefix(c.Text, "//rule:") {
			lines = append(lines, c.Text)
		}
	}
	return lines
}
