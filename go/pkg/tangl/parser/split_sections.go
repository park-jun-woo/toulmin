//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what splitSections — split file lines into tangl: section blocks
package parser

// splitSections scans lines for `## tangl:<Name>` headings and returns one
// section per heading, with content spanning to the next heading or EOF.
// Non-tangl content (titles, prose, other headings) is not included in any
// section, matching the spec's "파서 무시" rule for non-tangl text.
func splitSections(lines []string) []section {
	var secs []section
	i := 0
	for i < len(lines) {
		sec, next, ok := extractSection(lines, i)
		if !ok {
			i = next
			continue
		}
		secs = append(secs, sec)
		i = next
	}
	return secs
}
