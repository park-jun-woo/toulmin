//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what TestScanItemsParity — tests scanner item counts match the parser corpus statement counts
package scan

import (
	"os"
	"path/filepath"
	"testing"
)

// TestScanItemsParity locks the scanner's markdown recognition against the
// parser testdata. The parser's item tree is unexported, so the expected counts
// are the non-struck list-line totals of each fixture (verified against the
// files); Phase 002's encoder performs the live AST-vs-scanner dual check.
func TestScanItemsParity(t *testing.T) {
	var countItems func(items []Item) int
	countItems = func(items []Item) int {
		n := 0
		for _, it := range items {
			if !it.Struck {
				n++
			}
			n += countItems(it.Children)
		}
		return n
	}

	cases := map[string]int{"americano": 28, "access": 16, "transfer": 22}
	for name, want := range cases {
		src, err := os.ReadFile(filepath.Join("..", "..", "parser", "testdata", name+".md"))
		if err != nil {
			t.Fatalf("read %s: %v", name, err)
		}
		lines := ScanLines(src)
		_, secs := ScanSections(lines)
		got := 0
		for _, s := range secs {
			items, _ := ScanItems(s.Lines)
			got += countItems(items)
		}
		if got != want {
			t.Errorf("%s: scanner found %d non-struck items, want %d", name, got, want)
		}
	}
}
