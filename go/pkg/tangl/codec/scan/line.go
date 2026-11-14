//ff:type feature=tangl type=parser
//ff:what Line — one raw source line with its byte span, newline-excluded text, and number
package scan

// Line is one source line as produced by ScanLines. Span covers the raw bytes
// including the trailing newline (and any preceding \r), so consecutive Line
// spans tile the whole source. Text is the newline-excluded view of the line
// (leading indentation retained), and No is the 1-based line number.
type Line struct {
	Span Span
	Text string
	No   int
}
