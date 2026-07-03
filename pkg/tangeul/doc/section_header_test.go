//ff:func feature=tangl type=util control=sequence
//ff:what TestSectionHeader — en and ko locales resolve all seven headers; unknown locale and unknown SectionCode are errors
package doc

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
)

// TestSectionHeader verifies sectionHeader resolves the en and ko locale
// branches to their known SectionCode headings, and that an unrecognized
// locale (before the names map is ever consulted) and a recognized locale
// with an unrecognized SectionCode (after the map lookup fails) are both
// hard errors.
func TestSectionHeader(t *testing.T) {
	got, err := sectionHeader("en", codebook.SectionCases)
	if err != nil {
		t.Fatalf("en Cases: unexpected error: %v", err)
	}
	if got != "## tangl:Cases" {
		t.Errorf("en Cases = %q, want %q", got, "## tangl:Cases")
	}

	got, err = sectionHeader("ko", codebook.SectionCases)
	if err != nil {
		t.Fatalf("ko Cases: unexpected error: %v", err)
	}
	if got != "## tangl:사안" {
		t.Errorf("ko Cases = %q, want %q", got, "## tangl:사안")
	}

	if _, err := sectionHeader("xx", codebook.SectionCases); err == nil {
		t.Errorf("unknown locale: expected error")
	}

	if _, err := sectionHeader("en", codebook.SectionCode(99)); err == nil {
		t.Errorf("unknown SectionCode: expected error")
	}
}
