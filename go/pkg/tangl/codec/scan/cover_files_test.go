//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what TestScanFilesCover — tests full line tiling and lossless section partition over the real corpus
package scan

import (
	"os"
	"path/filepath"
	"testing"
)

// TestScanFilesCover sweeps the parser testdata and the test/ examples: every
// file's line table must tile [0, len) exactly, and within each section the
// item spans (recursively) plus prose spans must account for every content byte
// with no gap or overlap.
func TestScanFilesCover(t *testing.T) {
	var itemBytes func(items []Item) uint32
	itemBytes = func(items []Item) uint32 {
		var n uint32
		for _, it := range items {
			n += it.Span.Len + itemBytes(it.Children)
		}
		return n
	}
	sumLines := func(lns []Line) uint32 {
		var n uint32
		for _, ln := range lns {
			n += ln.Span.Len
		}
		return n
	}
	checkFile := func(p string) {
		src, err := os.ReadFile(p)
		if err != nil {
			t.Fatalf("read %s: %v", p, err)
		}
		lines := ScanLines(src)
		if err := VerifyCover(lines, len(src)); err != nil {
			t.Errorf("%s: line tiling incomplete: %v", p, err)
		}
		_, secs := ScanSections(lines)
		for _, s := range secs {
			items, prose := ScanItems(s.Lines)
			if itemBytes(items)+sumLines(prose) != sumLines(s.Lines) {
				t.Errorf("%s section %q: items+prose bytes != section bytes", p, s.Name)
			}
		}
	}

	paths := []string{
		filepath.Join("..", "..", "parser", "testdata", "americano.md"),
		filepath.Join("..", "..", "parser", "testdata", "access.md"),
		filepath.Join("..", "..", "parser", "testdata", "transfer.md"),
		filepath.Join("..", "..", "..", "..", "..", "test", "library-renewal.md"),
		filepath.Join("..", "..", "..", "..", "..", "test", "newsletter-tick-warn.md"),
		filepath.Join("..", "..", "..", "..", "..", "test", "newsletter-tick-fixed.md"),
		filepath.Join("..", "..", "..", "..", "..", "test", "order-saga.md"),
		filepath.Join("..", "..", "..", "..", "..", "test", "factory-line-cascade.md"),
	}
	for _, p := range paths {
		checkFile(p)
	}
}
