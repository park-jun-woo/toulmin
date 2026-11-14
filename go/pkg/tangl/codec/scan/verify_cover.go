//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what VerifyCover — assert a line table's spans tile [0, total) with no gaps or overlaps
package scan

import "fmt"

// VerifyCover checks that the spans of lines, taken in order, tile [0, total)
// exactly: each span must start where the previous ended, and the last must end
// at total. It is the scanner's self-check that ScanLines produced a complete
// raw-byte tiling, which is the precondition doc.VerifyTiling relies on.
func VerifyCover(lines []Line, total int) error {
	var next uint32
	for _, ln := range lines {
		if ln.Span.Off != next {
			return fmt.Errorf("line %d: span starts at %d, want %d (gap or overlap)", ln.No, ln.Span.Off, next)
		}
		next += ln.Span.Len
	}
	if int(next) != total {
		return fmt.Errorf("spans cover %d bytes, want %d", next, total)
	}
	return nil
}
