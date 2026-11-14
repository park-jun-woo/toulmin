//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseUndoSyntaxError — an undo edge missing its "when" clause is a parser error
package parser

import "testing"

// TestParseUndoSyntaxError verifies that "undo `f`" (missing "when `node`")
// is rejected as a malformed undo edge.
func TestParseUndoSyntaxError(t *testing.T) {
	src := "## tangl:Subject\n- this document is `t`\n\n## tangl:Cases\n- in case of `x`\n  - `a` is a general rule\n  - undo `f`\n"
	_, err := ParseSource(src, "undo_syntax_error.md")
	if err == nil {
		t.Fatal("expected error for malformed undo edge, got nil")
	}
}
