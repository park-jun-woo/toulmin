//ff:func feature=tangl type=parser control=sequence
//ff:what extractSection — build one section starting at a tangl: heading line
package parser

// extractSection builds the section starting at lines[i] if it is a
// `## tangl:<Name>` heading, along with the index to resume scanning from.
// ok is false when lines[i] is not a heading, in which case next is i+1.
func extractSection(lines []string, i int) (sec section, next int, ok bool) {
	name, isHeading := tanglSectionName(lines[i])
	if !isHeading {
		return section{}, i + 1, false
	}
	headerLine := i + 1
	start := i + 1
	end := findSectionEnd(lines, start)
	return section{
		Name:       name,
		HeaderLine: headerLine,
		Lines:      lines[start:end],
		LineOffset: start + 1,
	}, end, true
}
