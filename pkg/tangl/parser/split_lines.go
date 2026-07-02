//ff:func feature=tangl type=parser control=sequence
//ff:what splitLines — split source text into lines, normalizing CRLF
package parser

import "strings"

// splitLines splits src into lines, normalizing CRLF to LF first so line
// indices (1-based) match what an editor would show.
func splitLines(src string) []string {
	normalized := strings.ReplaceAll(src, "\r\n", "\n")
	return strings.Split(normalized, "\n")
}
