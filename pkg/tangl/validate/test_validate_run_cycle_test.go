//ff:func feature=tangl type=validator control=sequence
//ff:what TestValidateRunCycle — mutually run-cascading cases form a forbidden cycle
package validate

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/parser"
)

// TestValidateRunCycle verifies that case `x` running `y` while `y` runs
// `x` back is rejected as a run cascade cycle.
func TestValidateRunCycle(t *testing.T) {
	src := "## tangl:Subject\n- this document is `t`\n\n" +
		"## tangl:Cases\n" +
		"- in case of `x`\n" +
		"  - `a` is a general rule\n" +
		"  - run `y` when `a`\n\n" +
		"- in case of `y`\n" +
		"  - `b` is a general rule\n" +
		"  - run `x` when `b`\n"
	doc, err := parser.ParseSource(src, "run_cycle.md")
	if err != nil {
		t.Fatalf("ParseSource: %v", err)
	}
	err = Validate(doc)
	if err == nil {
		t.Fatal("expected run cascade cycle violation, got nil")
	}
	if !strings.Contains(err.Error(), "run cycle") {
		t.Errorf("error = %v, want mention of 'run cycle'", err)
	}
}
