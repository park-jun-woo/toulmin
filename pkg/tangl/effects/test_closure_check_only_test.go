//ff:func feature=tangl type=analyzer control=sequence
//ff:what TestClosureCheckOnly — a check-only endpoint always yields an empty closure
package effects

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/parser"
)

// TestClosureCheckOnly verifies that an endpoint with only a `check` clause
// (Evaluate mode, no execution) always yields an empty effect summary, even
// though the checked case has a `do` edge.
func TestClosureCheckOnly(t *testing.T) {
	src := "## tangl:Subject\n- this document is `t`\n\n" +
		"## tangl:See\n- see `f` from `pkg/f`\n\n" +
		"## tangl:Cases\n" +
		"- in case of `x`\n" +
		"  - `a` is a general rule\n" +
		"  - do `f`.`act` when `a`\n\n" +
		"## tangl:Provides\n" +
		"- provides `peek`\n" +
		"  - check `x`\n"
	doc, err := parser.ParseSource(src, "check_only.md")
	if err != nil {
		t.Fatalf("ParseSource: %v", err)
	}
	entries, err := Closure(doc, "peek")
	if err != nil {
		t.Fatalf("Closure: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("Closure(peek) = %+v, want empty slice", entries)
	}
}
