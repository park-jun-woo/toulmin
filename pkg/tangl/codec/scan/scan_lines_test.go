//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what TestScanLines — tests ScanLines for CRLF/LF mix, no-newline EOF, empty input, and tiling
package scan

import "testing"

func TestScanLines(t *testing.T) {
	// Mixed CRLF and LF, a blank line, and a final line without a newline.
	src := []byte("# Title\r\n\nlist\n- a")
	lines := ScanLines(src)
	type want struct {
		text string
		off  uint32
		ln   uint32
		no   int
	}
	wants := []want{
		{"# Title", 0, 9, 1}, // "# Title\r\n"
		{"", 9, 1, 2},        // "\n"
		{"list", 10, 5, 3},   // "list\n"
		{"- a", 15, 3, 4},    // "- a" (no trailing newline)
	}
	if len(lines) != len(wants) {
		t.Fatalf("got %d lines, want %d: %+v", len(lines), len(wants), lines)
	}
	for i, w := range wants {
		got := lines[i]
		if got.Text != w.text || got.Span.Off != w.off || got.Span.Len != w.ln || got.No != w.no {
			t.Errorf("line %d = %+v, want text=%q off=%d len=%d no=%d", i, got, w.text, w.off, w.ln, w.no)
		}
	}
	if err := VerifyCover(lines, len(src)); err != nil {
		t.Errorf("VerifyCover: %v", err)
	}

	// Empty input yields no lines and trivially tiles [0, 0).
	empty := ScanLines([]byte(""))
	if len(empty) != 0 {
		t.Errorf("empty src: got %d lines, want 0", len(empty))
	}
	if err := VerifyCover(empty, 0); err != nil {
		t.Errorf("VerifyCover(empty): %v", err)
	}
}
