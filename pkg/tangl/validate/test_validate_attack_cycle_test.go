//ff:func feature=tangl type=validator control=sequence
//ff:what TestValidateAttackCycle — mutually attacking nodes form a forbidden defeat cycle
package validate

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/parser"
)

// TestValidateAttackCycle verifies that node `a` attacking `b` while `b`
// attacks `a` back, within the same case, is rejected as a don't-edge cycle.
func TestValidateAttackCycle(t *testing.T) {
	src := "## tangl:Subject\n- this document is `t`\n\n" +
		"## tangl:Cases\n" +
		"- in case of `x`\n" +
		"  - `a` is a general rule\n" +
		"  - `b` is a general rule\n" +
		"  - don't `a` when `b`\n" +
		"  - don't `b` when `a`\n"
	doc, err := parser.ParseSource(src, "attack_cycle.md")
	if err != nil {
		t.Fatalf("ParseSource: %v", err)
	}
	err = Validate(doc)
	if err == nil {
		t.Fatal("expected don't-edge cycle violation, got nil")
	}
	if !strings.Contains(err.Error(), "don't cycle") {
		t.Errorf("error = %v, want mention of \"don't cycle\"", err)
	}
}
