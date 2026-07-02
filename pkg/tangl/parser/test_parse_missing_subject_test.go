//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseMissingSubject — a document without tangl:Subject is a parser error
package parser

import "testing"

// TestParseMissingSubject verifies that a document with tangl:Cases but no
// tangl:Subject section is rejected (Subject is required by the parser, not
// deferred to the validator).
func TestParseMissingSubject(t *testing.T) {
	src := "## tangl:Cases\n- in case of `x`\n  - `a` is required\n"
	_, err := ParseSource(src, "missing_subject.md")
	if err == nil {
		t.Fatal("expected error for missing tangl:Subject, got nil")
	}
}
