//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what countIndent — count leading whitespace characters in a line
package parser

// countIndent counts leading whitespace characters, treating tabs as 4 spaces.
func countIndent(line string) int {
	count := 0
	for _, ch := range line {
		if ch == ' ' {
			count++
		} else if ch == '\t' {
			count += 4
		} else {
			break
		}
	}
	return count
}
