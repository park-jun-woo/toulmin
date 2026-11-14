//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what TestScanSections — tests ScanSections for preamble, scope-closing headings, and verbatim names
package scan

import "testing"

func TestScanSections(t *testing.T) {
	src := []byte("intro prose\n- stray list\n\n" +
		"## tangl:Subject\n- this is `x`\n\n" +
		"## Other Heading\nprose under other\n\n" +
		"## tangl:주제\n- 이것은 `y`\n")
	lines := ScanLines(src)
	preamble, secs := ScanSections(lines)

	// Preamble is only the lines before the first tangl section.
	if len(preamble) != 3 {
		t.Fatalf("preamble = %d lines, want 3: %+v", len(preamble), preamble)
	}
	if preamble[0].Text != "intro prose" || preamble[1].Text != "- stray list" {
		t.Errorf("preamble content = %q / %q", preamble[0].Text, preamble[1].Text)
	}

	if len(secs) != 2 {
		t.Fatalf("got %d sections, want 2: %+v", len(secs), secs)
	}

	// A non-tangl heading closes the Subject section's scope.
	if secs[0].Name != "Subject" {
		t.Errorf("secs[0].Name = %q, want Subject", secs[0].Name)
	}
	if len(secs[0].Lines) != 2 { // "- this is `x`" and the blank line
		t.Errorf("secs[0].Lines = %d, want 2: %+v", len(secs[0].Lines), secs[0].Lines)
	}
	for _, ln := range secs[0].Lines {
		if ln.Text == "prose under other" {
			t.Errorf("Subject section leaked prose past the closing heading")
		}
	}

	// The Korean section name is preserved verbatim (no canonicalization here).
	if secs[1].Name != "주제" {
		t.Errorf("secs[1].Name = %q, want 주제", secs[1].Name)
	}
	// HeaderSpan points at the heading line's raw bytes.
	hs := secs[1].HeaderSpan
	if string(src[hs.Off:hs.Off+hs.Len]) != "## tangl:주제\n" {
		t.Errorf("secs[1].HeaderSpan bytes = %q", string(src[hs.Off:hs.Off+hs.Len]))
	}
}
