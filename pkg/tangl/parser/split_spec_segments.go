//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what splitSpecSegments — split spec call text by comma and 'and' delimiters respecting parentheses
package parser

// splitSpecSegments splits spec call text by commas and "and" keywords, respecting parenthesis depth.
func splitSpecSegments(text string) []string {
	var segments []string
	depth := 0
	start := 0
	for i := 0; i < len(text); i++ {
		ch := text[i]
		if ch == '(' {
			depth++
		} else if ch == ')' {
			depth--
		}
		if depth != 0 {
			continue
		}
		seg, newStart, advanced := trySegmentSplit(text, i, start)
		if !advanced {
			continue
		}
		segments = append(segments, seg)
		start = newStart
		i = start - 1
	}
	if start < len(text) {
		segments = append(segments, text[start:])
	}
	return segments
}
