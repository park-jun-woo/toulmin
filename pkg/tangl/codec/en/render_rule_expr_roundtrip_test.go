//ff:func feature=tangl type=parser control=sequence
//ff:what TestRenderRuleExprRoundtrip — round-trip new-syntax Rules (one-line heads, nested exists, operands) through the pipeline
package en

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/pkg/tangl/parser"
)

// TestRenderRuleExprRoundtrip mirrors TestRenderRuleRoundtrip for the expression
// grammar: it authors a rule set whose top-level conditions are a lone predicate
// Call and a lone Exists written one-line (so encode -> render reproduces the same
// minimal one-line form, pinning the (a)-decision idempotence), plus nesting the
// fixture does not cover — an exists whose body is an and/or chain, a two-tier
// nested exists (inner condition references the outer binder), an operand compare
// inside a block "not", and an exists term inside an "either" group. Judgment is
// encode -> doc.Render -> re-parse -> EqualIgnoringLine.
func TestRenderRuleExprRoundtrip(t *testing.T) {
	src := "## tangl:Subject\n- this document is `demo`\n\n" +
		"## tangl:See\n- see `check` from `mail/check`\n\n" +
		"## tangl:Rules\n" +
		"- `spammy` when `check`.`isSpam` holds\n" +
		"- `hasCoupon` when any of `coupons` as `c` satisfies `c`.`minPrice` is at most `basePrice`\n" +
		"- `activePromotion` when any of `promotions` as `p` satisfies\n" +
		"  - `p`.`active` equals true\n" +
		"  - or `p`.`discount` is at least 10\n" +
		"- `nestedExists` when any of `orders` as `o` satisfies\n" +
		"  - any of `o`.`items` as `i` satisfies `i`.`price` is at most `o`.`budget`\n" +
		"- `notOperand` when\n" +
		"  - not `author`.`trustScore` is at least 0.7\n" +
		"- `eitherExists` when\n" +
		"  - either\n" +
		"    - any of `coupons` as `c` satisfies `c`.`minPrice` is at most `basePrice`\n" +
		"    - and `credit score` is at least 700\n\n" +
		"## tangl:Cases\n- in case of `x`\n  - `spammy` is a general rule\n"

	d, err := ENCodec{}.Encode([]byte(src), "rule_expr.md")
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
