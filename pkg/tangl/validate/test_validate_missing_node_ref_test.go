//ff:func feature=tangl type=validator control=sequence
//ff:what TestValidateMissingNodeRef — a don't edge naming an unregistered attacker is a violation
package validate

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/parser"
)

// TestValidateMissingNodeRef verifies that `don't 'a' when 'ghost'`, where
// `ghost` is never registered as a node, is reported by Validate.
func TestValidateMissingNodeRef(t *testing.T) {
	src := "## tangl:Subject\n- this document is `t`\n\n" +
		"## tangl:Cases\n" +
		"- in case of `x`\n" +
		"  - `a` is a general rule\n" +
		"  - don't `a` when `ghost`\n"
	doc, err := parser.ParseSource(src, "missing_ref.md")
	if err != nil {
		t.Fatalf("ParseSource: %v", err)
	}
	err = Validate(doc)
	if err == nil {
		t.Fatal("expected missing node reference violation, got nil")
	}
	if !strings.Contains(err.Error(), `"ghost"`) {
		t.Errorf("error = %v, want mention of \"ghost\"", err)
	}
}
