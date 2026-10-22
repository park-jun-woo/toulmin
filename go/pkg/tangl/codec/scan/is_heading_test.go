//ff:func feature=tangl type=parser control=sequence
//ff:what TestIsHeading — tests isHeading for heading prefixes, leading whitespace, and non-heading text
package scan

import "testing"

func TestIsHeading(t *testing.T) {
	if !isHeading("# Title") {
		t.Errorf("isHeading(%q) = false, want true", "# Title")
	}
	if !isHeading("## tangl:Subject") {
		t.Errorf("isHeading(%q) = false, want true", "## tangl:Subject")
	}
	// Leading whitespace is trimmed before the prefix check.
	if !isHeading("  # indented") {
		t.Errorf("isHeading(%q) = false, want true", "  # indented")
	}
	if isHeading("prose line") {
		t.Errorf("isHeading(%q) = true, want false", "prose line")
	}
	if isHeading("") {
		t.Errorf("isHeading(%q) = true, want false", "")
	}
	if isHeading("  ") {
		t.Errorf("isHeading(%q) = true, want false", "  ")
	}
}
