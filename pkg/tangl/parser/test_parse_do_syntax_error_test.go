//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseDoSyntaxError — a do edge missing its "when" clause is a parser error
package parser

import "testing"

// TestParseDoSyntaxError verifies that "do `f` if" (missing "when `node`")
// is rejected as a malformed do edge.
func TestParseDoSyntaxError(t *testing.T) {
	src := "## tangl:Subject\n- this document is `t`\n\n## tangl:Cases\n- in case of `x`\n  - `a` is a general rule\n  - do `f` if\n"
	_, err := ParseSource(src, "do_syntax_error.md")
	if err == nil {
		t.Fatal("expected error for malformed do edge, got nil")
	}
}
