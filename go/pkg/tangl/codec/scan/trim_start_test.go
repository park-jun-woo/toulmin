//ff:func feature=tangl type=parser control=sequence
//ff:what TestTrimStart — tests trimStart for leading spaces, tabs, no leading whitespace, all-whitespace, and empty input
package scan

import "testing"

func TestTrimStart(t *testing.T) {
	if got := trimStart("  x"); got != "x" {
		t.Errorf("trimStart(spaces) = %q, want %q", got, "x")
	}
	if got := trimStart("\t\tx"); got != "x" {
		t.Errorf("trimStart(tabs) = %q, want %q", got, "x")
	}
	if got := trimStart(" \t x"); got != "x" {
		t.Errorf("trimStart(mixed) = %q, want %q", got, "x")
	}
	// No leading whitespace: the loop body never runs.
	if got := trimStart("x"); got != "x" {
		t.Errorf("trimStart(no leading ws) = %q, want %q", got, "x")
	}
	// All whitespace: the loop runs until i == len(s).
	if got := trimStart("   "); got != "" {
		t.Errorf("trimStart(all ws) = %q, want %q", got, "")
	}
	// Empty input: the loop body never runs.
	if got := trimStart(""); got != "" {
		t.Errorf("trimStart(empty) = %q, want %q", got, "")
	}
}
