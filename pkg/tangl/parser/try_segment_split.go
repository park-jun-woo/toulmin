//ff:func feature=tangl type=parser control=sequence
//ff:what trySegmentSplit — try to split a segment at the current position by 'and' or comma delimiter
package parser

import "strings"

// trySegmentSplit checks if a segment boundary exists at position i and returns the segment and new start.
func trySegmentSplit(text string, i int, start int) (string, int, bool) {
	rest := strings.TrimSpace(text[i+1:])
	if strings.HasPrefix(rest, "and ") {
		seg := text[start : i+1]
		newStart := skipSpaces(text, i+1)
		if newStart+4 <= len(text) && text[newStart:newStart+4] == "and " {
			newStart += 4
		}
		return seg, newStart, true
	}
	if strings.HasPrefix(rest, ", ") || strings.HasPrefix(rest, ",") {
		seg := text[start : i+1]
		newStart := skipCommaSpaces(text, i+1)
		return seg, newStart, true
	}
	return "", 0, false
}
