//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what ScanLines — build a raw-byte line offset table over the source
package scan

// ScanLines splits src into lines on each `\n`, keeping every line's span over
// the raw bytes: the span includes the trailing `\n` (and its preceding `\r`),
// so concatenating all spans tiles [0, len(src)) with no gaps or overlaps. Text
// excludes that newline. A final line without a trailing newline runs to EOF;
// empty input yields no lines.
func ScanLines(src []byte) []Line {
	var lines []Line
	off := 0
	no := 1
	for off < len(src) {
		end := off
		for end < len(src) && src[end] != '\n' {
			end++
		}
		spanEnd := end
		if spanEnd < len(src) {
			spanEnd++ // include the '\n'
		}
		textEnd := end
		if textEnd > off && src[textEnd-1] == '\r' {
			textEnd-- // drop the '\r' of a CRLF pair from Text
		}
		lines = append(lines, Line{
			Span: Span{Off: uint32(off), Len: uint32(spanEnd - off)},
			Text: string(src[off:textEnd]),
			No:   no,
		})
		off = spanEnd
		no++
	}
	return lines
}
