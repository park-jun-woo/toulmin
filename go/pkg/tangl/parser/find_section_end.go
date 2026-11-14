//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what findSectionEnd — find the index of the next heading line or EOF
package parser

// findSectionEnd scans forward from start and returns the index of the next
// heading line (any level), or len(lines) if none remains.
func findSectionEnd(lines []string, start int) int {
	i := start
	for i < len(lines) && !isHeadingLine(lines[i]) {
		i++
	}
	return i
}
