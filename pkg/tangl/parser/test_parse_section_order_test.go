//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseSectionOrder — sections out of Subject→...→Internal order is an error
package parser

import "testing"

// TestParseSectionOrder verifies that tangl:Cases appearing before
// tangl:Subject (violating the required section order) is rejected.
func TestParseSectionOrder(t *testing.T) {
	src := "## tangl:Cases\n- in case of `x`\n  - `a` is required\n\n## tangl:Subject\n- this document is `t`\n"
	_, err := ParseSource(src, "section_order.md")
	if err == nil {
		t.Fatal("expected error for out-of-order sections, got nil")
	}
}
