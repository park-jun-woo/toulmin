//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what sectionEnd — find the index of the next heading line or EOF from start
package scan

// sectionEnd scans lines forward from start and returns the index of the next
// heading line (any level), or len(lines) if none remains. It mirrors the
// parser's findSectionEnd so a section's content stops at the same boundary.
func sectionEnd(lines []Line, start int) int {
	i := start
	for i < len(lines) && !isHeading(lines[i].Text) {
		i++
	}
	return i
}
