//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what rulesExprByName — parse the extended-expression fixture and index its rules by name
package parser

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// rulesExprByName parses testdata/rules_expr.md and returns its inline rules
// indexed by name, for the golden assertions below.
func rulesExprByName(t *testing.T) map[string]ast.InlineRule {
	t.Helper()
	doc, err := Parse("testdata/rules_expr.md")
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	rules := map[string]ast.InlineRule{}
	for _, r := range doc.Rules {
		rules[r.Name] = r
	}
	return rules
}
