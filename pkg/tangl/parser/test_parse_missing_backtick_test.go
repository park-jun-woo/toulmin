//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseMissingBacktick — subject name without backticks is a parser error
package parser

import "testing"

// TestParseMissingBacktick verifies that a user-defined name written without
// backticks (a bare keyword-looking token) is rejected.
func TestParseMissingBacktick(t *testing.T) {
	src := "## tangl:Subject\n- this document is americano\n\n## tangl:Cases\n- in case of `x`\n  - `a` is required\n"
	_, err := ParseSource(src, "missing_backtick.md")
	if err == nil {
		t.Fatal("expected error for missing backtick, got nil")
	}
}
