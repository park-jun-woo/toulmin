//ff:func feature=tangl type=parser control=sequence
//ff:what TestRenderStmtExtra — round-trip statement forms absent from the fixture corpus
package en

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangeul"
	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/parser"
)

// TestRenderStmtExtra covers statement forms the fixtures do not exercise — struct
// definitions with `has` fields, a node with `with`+`qualified`, `checking`,
// require-as-type, a do-once-certainty edge, and an `on <event>` internal — by
// encoding a synthetic document, rendering it (mode 2), re-parsing, and asserting
// the AST is unchanged up to Line.
func TestRenderStmtExtra(t *testing.T) {
	src := "## tangl:Subject\n- this document is `demo`\n\n" +
		"## tangl:See\n- see `p` from `pkg/p`\n\n" +
		"## tangl:Definitions\n" +
		"- `threshold` means 10 as `p`.`T`\n" +
		"- `applicant` means\n" +
		"  - has `income` as Currency\n" +
		"  - has `spec` as `p`.`Spec`\n\n" +
		"## tangl:Cases\n" +
		"- in case of `can go`\n" +
		"  - `user` is required\n" +
		"  - `age` is required as Integer\n" +
		"  - `ok` is a general rule using `p`.`check` with `threshold` qualified 0.75\n" +
		"  - `bad` is a counter rule using `p`.`isBad`\n" +
		"  - `warn` is an except rule using `p`.`warn`\n" +
		"  - don't `ok` when `bad`\n" +
		"  - do `p`.`act` once when `ok` if at least 80% certain\n" +
		"  - undo `p`.`revert` when `ok`\n" +
		"  - run `can sub` when `ok`\n" +
		"- in case of `can sub`\n" +
		"  - `x` is required\n" +
		"  - `y` is a general rule checking `can go`\n" +
		"  - do `p`.`sub` when `y`\n\n" +
		"## tangl:Provides\n" +
		"- provides `entry`\n" +
		"  - `user` is required\n" +
		"  - run `can go`\n" +
		"  - check `can sub`\n\n" +
		"## tangl:Internal\n" +
		"- on `p`.`event`\n" +
		"  - run `can go`\n" +
		"- every 30s until `can go`\n" +
		"  - run `can go`\n"

	d, err := ENCodec{}.Encode([]byte(src), "extra.md")
	if err != nil {
		t.Fatalf("Encode: %v", err)
	}
	rendered, err := doc.Render(d, ENCodec{})
	if err != nil {
		t.Fatalf("Render: %v", err)
	}
	reparsed, err := parser.ParseSource(string(rendered), d.AST.Path)
	if err != nil {
		t.Fatalf("re-parse failed: %v\n---\n%s", err, rendered)
	}
	if !tangeul.EqualIgnoringLine(d.AST, reparsed) {
		t.Errorf("re-parsed AST differs\n--- rendered ---\n%s", rendered)
	}
}
