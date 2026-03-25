//ff:func feature=tangl type=parser control=sequence
//ff:what parseSectionEntry — parse a single entry within a TANGL section and append to File
package parser

import "fmt"

// parseSectionEntry processes one line in a section. Returns lines consumed and updated errors.
func parseSectionEntry(lines []string, i int, sectionName string, trimmedRaw string, f *File, errs []string) (int, []string) {
	if trimmedRaw == "" {
		return 1, errs
	}

	stripped := stripListPrefix(trimmedRaw)
	if stripped == "" {
		return 1, errs
	}

	parentGraph := ""
	indent := countIndent(lines[i])
	if indent >= 2 {
		parentGraph = findParentGraph(lines, i, indent)
	}

	lineNum := i + 1
	consumed, node, err := parseSectionLine(lines, i, sectionName, stripped, lineNum, parentGraph)
	if err != nil {
		errs = append(errs, fmt.Sprintf("line %d: %s", lineNum, err))
		return consumed, errs
	}
	appendNode(f, node)
	return consumed, errs
}
