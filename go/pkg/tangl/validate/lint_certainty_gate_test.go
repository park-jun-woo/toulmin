//ff:func feature=tangl type=validator control=sequence
//ff:what TestLintCertaintyGate — Lint surfaces an unreachable certainty gate from a parsed document
package validate

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/parser"
)

func TestLintCertaintyGate(t *testing.T) {
	src := "## tangl:Subject\n- this document is `g`\n\n" +
		"## tangl:Cases\n- in case of `c`\n  - `n` is a general rule\n" +
		"  - do `f` when `n` if at most 50% certain\n"
	doc, err := parser.ParseSource(src, "g.md")
	if err != nil {
		t.Fatalf("ParseSource: %v", err)
	}
	w := Lint(doc)
	if len(w) != 1 || !strings.Contains(w[0], "never passes") {
		t.Fatalf("Lint = %v, want one never-passes warning", w)
	}
}
