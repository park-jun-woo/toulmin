//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what entriesFromLines — reduce raw lines to non-blank indent/text entries
package parser

import "strings"

// entriesFromLines drops blank lines and records each remaining line's
// indent width, trimmed text, and 1-based file line number.
func entriesFromLines(lines []string, lineOffset int) []lineEntry {
	var entries []lineEntry
	for i, raw := range lines {
		if strings.TrimSpace(raw) == "" {
			continue
		}
		entries = append(entries, lineEntry{
			Indent: lineIndent(raw),
			Text:   strings.TrimSpace(raw),
			Line:   lineOffset + i,
		})
	}
	return entries
}
