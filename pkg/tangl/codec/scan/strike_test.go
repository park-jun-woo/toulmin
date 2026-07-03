//ff:func feature=tangl type=parser control=sequence
//ff:what TestStrike — tests strike for wrapped, unwrapped, prefix-only, and too-short overlapping content
package scan

import "testing"

func TestStrike(t *testing.T) {
	// Wrapped content is unwrapped and flagged.
	if text, struck := strike("~~gone~~"); !struck || text != "gone" {
		t.Errorf("strike(wrapped) = (%q, %v), want (\"gone\", true)", text, struck)
	}

	// No "~~" prefix at all.
	if text, struck := strike("plain"); struck || text != "plain" {
		t.Errorf("strike(plain) = (%q, %v), want (\"plain\", false)", text, struck)
	}

	// Prefix present but no matching suffix.
	if text, struck := strike("~~onlyprefix"); struck || text != "~~onlyprefix" {
		t.Errorf("strike(prefix only) = (%q, %v), want unstruck passthrough", text, struck)
	}

	// Exactly four characters: prefix and suffix overlap fully, content is empty.
	if text, struck := strike("~~~~"); !struck || text != "" {
		t.Errorf("strike(\"~~~~\") = (%q, %v), want (\"\", true)", text, struck)
	}

	// Prefix and suffix both match (overlapping) but length is under 4.
	if text, struck := strike("~~~"); struck || text != "~~~" {
		t.Errorf("strike(\"~~~\") = (%q, %v), want unstruck passthrough", text, struck)
	}
}
