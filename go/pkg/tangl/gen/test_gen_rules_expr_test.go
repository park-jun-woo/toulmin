//ff:func feature=tangl type=codegen control=iteration dimension=1
//ff:what TestGenRulesExpr — Generate compiles the extended-expression fixture into gofmt-clean closure code
package gen

import (
	"go/format"
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/parser"
)

// TestGenRulesExpr generates the extended-expression fixture (field paths,
// spec refs, host value/predicate calls, exists quantifiers) and checks the
// output re-formats cleanly and contains the closure-composition patterns:
// the gated tangl* helpers, the alias-qualified host call sites, and the See
// imports those call sites need.
func TestGenRulesExpr(t *testing.T) {
	doc, err := parser.Parse("../parser/testdata/rules_expr.md")
	if err != nil {
		t.Fatalf("Parse(rules_expr.md): %v", err)
	}
	src, err := Generate(doc)
	if err != nil {
		t.Fatalf("Generate(rules_expr.md): %v", err)
	}
	if _, err := format.Source([]byte(src)); err != nil {
		t.Fatalf("re-format.Source: %v\n--- source ---\n%s", err, src)
	}
	for _, want := range []string{
		"func tanglAny(",
		"func tanglHolds(",
		"func tanglSpecVal(",
		"func tanglCtxVal(",
		"func tanglWalk(",
		"func tanglCompareVal(",
		"func tanglOr(",
		"func tanglAnd(",
		"nlp.score(ctx, specs)",
		"check.isSpam(ctx, specs)",
		`tanglSpecVal(specs, "discount", "min price")`,
		`tanglCtxVal(ctx, "author", "trustScore")`,
		`"coupons[" + strconv.Itoa(i0) + "].minPrice"`,
		`"any of coupons as c satisfies"`,
		`tanglCompareVal(e0, true, "equals", rv)`,
		`nlp "`,
		`check "`,
	} {
		if !strings.Contains(src, want) {
			t.Errorf("generated source missing %q\n--- source ---\n%s", want, src)
		}
	}
}
