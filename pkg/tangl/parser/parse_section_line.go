//ff:func feature=tangl type=parser control=sequence
//ff:what parseSectionLine — parse a single line within a section, handling inline rules specially
package parser

import "strings"

// parseSectionLine parses a single line within a section. Returns lines consumed, the node, and any error.
func parseSectionLine(lines []string, i int, sectionName string, stripped string, lineNum int, parentGraph string) (int, any, error) {
	if sectionName == "Rules" && strings.HasPrefix(stripped, "rule \"") && strings.HasSuffix(stripped, "\" is") {
		body := ""
		if i+1 < len(lines) {
			body = strings.TrimSpace(lines[i+1])
		}
		node, err := parseInlineRule(stripped, body, lineNum)
		if err != nil {
			return 2, nil, err
		}
		return 2, node, nil
	}

	node, err := parseLine(stripped, lineNum, parentGraph)
	if err != nil {
		return 1, nil, err
	}
	return 1, node, nil
}
