//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what VerifyTiling — check Blocks tile [0, len(Source)) in stream order with no gap/overlap
package doc

import (
	"errors"
	"fmt"
)

// VerifyTiling checks that Blocks, concatenated in stream order, tile
// [0, len(Source)) with no gap, overlap, or order reversal (README contract §4,
// the G3 expansion gate). The error names the first violation: its offset and
// block index. An empty Source with no Blocks is valid; leftover Source bytes are
// a trailing gap.
func VerifyTiling(d *Doc) error {
	if d == nil {
		return errors.New("tangeul/doc: VerifyTiling on nil Doc")
	}
	var expected uint32
	var prevOff uint32
	for i, b := range d.Blocks {
		off := b.Span.Off
		switch {
		case off > expected:
			return fmt.Errorf("tangeul/doc: tiling gap before block %d at offset %d (expected %d)", i, off, expected)
		case off < expected && i > 0 && off < prevOff:
			return fmt.Errorf("tangeul/doc: tiling order reversal at block %d: offset %d precedes previous block offset %d", i, off, prevOff)
		case off < expected:
			return fmt.Errorf("tangeul/doc: tiling overlap at block %d: offset %d overlaps previous end %d", i, off, expected)
		}
		prevOff = off
		expected = off + b.Span.Len
	}
	if expected != uint32(len(d.Source)) {
		return fmt.Errorf("tangeul/doc: tiling trailing gap: covered %d of %d Source bytes", expected, len(d.Source))
	}
	return nil
}
