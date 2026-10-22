//ff:func feature=tangl type=parser control=sequence
//ff:what TestSectionName — tests sectionName for missing prefix, empty name after the prefix, and a valid name
package scan

import "testing"

func TestSectionName(t *testing.T) {
	// No "## tangl:" prefix at all.
	if name, ok := sectionName("## Other Heading"); ok || name != "" {
		t.Errorf("sectionName(non-tangl) = (%q, %v), want (\"\", false)", name, ok)
	}

	// Prefix present but the name is blank (whitespace only).
	if name, ok := sectionName("## tangl:   "); ok || name != "" {
		t.Errorf("sectionName(blank name) = (%q, %v), want (\"\", false)", name, ok)
	}

	// Valid tangl section heading, with leading whitespace on the line itself.
	if name, ok := sectionName("  ## tangl:Subject  "); !ok || name != "Subject" {
		t.Errorf("sectionName(valid) = (%q, %v), want (\"Subject\", true)", name, ok)
	}

	// Non-ASCII names are preserved verbatim.
	if name, ok := sectionName("## tangl:주제"); !ok || name != "주제" {
		t.Errorf("sectionName(unicode) = (%q, %v), want (\"주제\", true)", name, ok)
	}
}
