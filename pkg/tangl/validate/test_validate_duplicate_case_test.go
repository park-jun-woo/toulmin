//ff:func feature=tangl type=validator control=sequence
//ff:what TestValidateDuplicateCase — two "in case of" blocks with the same name is a violation
package validate

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/parser"
)

// TestValidateDuplicateCase verifies that declaring `x` as a case name twice
// is reported by Validate.
func TestValidateDuplicateCase(t *testing.T) {
	src := "## tangl:Subject\n- this document is `t`\n\n" +
		"## tangl:Cases\n" +
		"- in case of `x`\n" +
		"  - `a` is a general rule\n\n" +
		"- in case of `x`\n" +
		"  - `b` is a general rule\n"
	doc, err := parser.ParseSource(src, "dup_case.md")
	if err != nil {
		t.Fatalf("ParseSource: %v", err)
	}
	err = Validate(doc)
	if err == nil {
		t.Fatal("expected duplicate case name violation, got nil")
	}
	if !strings.Contains(err.Error(), "duplicate case name") {
		t.Errorf("error = %v, want mention of 'duplicate case name'", err)
	}
}
