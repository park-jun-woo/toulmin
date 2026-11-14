//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseUnknownSection — an unrecognized tangl: section name is a parser error
package parser

import "testing"

// TestParseUnknownSection verifies that a `## tangl:Bogus` heading (not one
// of the seven recognized sections) is rejected.
func TestParseUnknownSection(t *testing.T) {
	src := "## tangl:Bogus\n- something\n"
	_, err := ParseSource(src, "unknown_section.md")
	if err == nil {
		t.Fatal("expected error for unknown tangl section, got nil")
	}
}
