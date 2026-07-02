//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseMissingCases — a document with tangl:Subject but no tangl:Cases is a parser error
package parser

import "testing"

// TestParseMissingCases verifies that a document with a tangl:Subject
// section but no tangl:Cases section is rejected (Cases is required by the
// parser, not deferred to the validator).
func TestParseMissingCases(t *testing.T) {
	src := "## tangl:Subject\n- this document is `t`\n"
	_, err := ParseSource(src, "missing_cases.md")
	if err == nil {
		t.Fatal("expected error for missing tangl:Cases, got nil")
	}
}
