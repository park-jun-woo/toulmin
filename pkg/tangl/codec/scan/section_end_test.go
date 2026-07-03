//ff:func feature=tangl type=parser control=sequence
//ff:what TestSectionEnd — tests sectionEnd for an immediate heading, no heading before EOF, and start already at len
package scan

import "testing"

func TestSectionEnd(t *testing.T) {
	lines := []Line{
		{Text: "content one", No: 1},
		{Text: "content two", No: 2},
		{Text: "## Next Heading", No: 3},
		{Text: "trailing content", No: 4},
	}

	// Stops at the heading a few lines in (loop runs, then condition fails
	// because isHeading becomes true).
	if got := sectionEnd(lines, 0); got != 2 {
		t.Errorf("sectionEnd(lines, 0) = %d, want 2", got)
	}

	// No heading remains before EOF: loop runs until i == len(lines).
	noHeading := []Line{
		{Text: "a", No: 1},
		{Text: "b", No: 2},
	}
	if got := sectionEnd(noHeading, 0); got != len(noHeading) {
		t.Errorf("sectionEnd(noHeading, 0) = %d, want %d", got, len(noHeading))
	}

	// start already at len(lines): loop body never runs.
	if got := sectionEnd(lines, len(lines)); got != len(lines) {
		t.Errorf("sectionEnd(lines, len) = %d, want %d", got, len(lines))
	}

	// The heading itself is at start: loop body never runs either.
	if got := sectionEnd(lines, 2); got != 2 {
		t.Errorf("sectionEnd(lines, 2) = %d, want 2", got)
	}
}
