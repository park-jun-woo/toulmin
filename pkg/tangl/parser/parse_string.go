//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what ParseString — parse TANGL markdown content into File AST
package parser

import (
	"fmt"
	"strings"
)

// ParseString parses TANGL markdown content and returns the File AST.
func ParseString(content string) (*File, error) {
	lines := strings.Split(content, "\n")
	f := &File{}
	var errs []string

	i := 0
	for i < len(lines) {
		line := lines[i]
		trimmed := strings.TrimSpace(line)

		if !strings.HasPrefix(trimmed, "## tangl:") {
			i++
			continue
		}

		sectionName := strings.TrimPrefix(trimmed, "## tangl:")
		i++
		i, errs = parseSection(lines, i, sectionName, f, errs)
	}

	if len(errs) > 0 {
		return nil, fmt.Errorf("parse errors:\n%s", strings.Join(errs, "\n"))
	}
	return f, nil
}
