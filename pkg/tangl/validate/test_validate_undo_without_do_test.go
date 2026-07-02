//ff:func feature=tangl type=validator control=sequence
//ff:what TestValidateUndoWithoutDo — an undo edge with no preceding do on its node is a violation
package validate

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/parser"
)

// TestValidateUndoWithoutDo verifies that `undo 'f' when 'a'` with no
// preceding `do ... when 'a'` in the same case is reported by Validate
// (spec §undo).
func TestValidateUndoWithoutDo(t *testing.T) {
	src := "## tangl:Subject\n- this document is `t`\n\n" +
		"## tangl:See\n- see `bank` from `bank/core`\n\n" +
		"## tangl:Cases\n" +
		"- in case of `x`\n" +
		"  - `a` is a general rule\n" +
		"  - undo `bank`.`f` when `a`\n"
	doc, err := parser.ParseSource(src, "undo_without_do.md")
	if err != nil {
		t.Fatalf("ParseSource: %v", err)
	}
	err = Validate(doc)
	if err == nil {
		t.Fatal("expected undo-without-do violation, got nil")
	}
	if !strings.Contains(err.Error(), "no preceding do") {
		t.Errorf("error = %v, want mention of 'no preceding do'", err)
	}
}
