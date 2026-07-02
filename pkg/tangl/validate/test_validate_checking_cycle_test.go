//ff:func feature=tangl type=validator control=sequence
//ff:what TestValidateCheckingCycle — mutually "checking" cases form a forbidden cycle
package validate

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/parser"
)

// TestValidateCheckingCycle verifies that case `x` checking `y` while `y`
// checks `x` back is rejected as a checking cycle.
func TestValidateCheckingCycle(t *testing.T) {
	src := "## tangl:Subject\n- this document is `t`\n\n" +
		"## tangl:Cases\n" +
		"- in case of `x`\n" +
		"  - `a` is a general rule checking `y`\n\n" +
		"- in case of `y`\n" +
		"  - `b` is a general rule checking `x`\n"
	doc, err := parser.ParseSource(src, "checking_cycle.md")
	if err != nil {
		t.Fatalf("ParseSource: %v", err)
	}
	err = Validate(doc)
	if err == nil {
		t.Fatal("expected checking cycle violation, got nil")
	}
	if !strings.Contains(err.Error(), "checking cycle") {
		t.Errorf("error = %v, want mention of 'checking cycle'", err)
	}
}
