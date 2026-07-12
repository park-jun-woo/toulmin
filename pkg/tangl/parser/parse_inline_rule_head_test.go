//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseInlineRuleHead — tests parseInlineRuleHead dispatch to cond-list, exists, call, and compare branches
package parser

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestParseInlineRuleHead(t *testing.T) {
	t.Run("EmptyRestCondList", func(t *testing.T) {
		it := item{
			Text:     "`ok` when",
			Line:     1,
			Children: []item{{Text: "`amount` equals 5", Line: 2}},
		}
		expr, err := parseInlineRuleHead("", it, "test.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if _, ok := expr.(ast.Compare); !ok {
			t.Fatalf("expected Compare from cond list, got %T", expr)
		}
	})

	t.Run("QuantifierExists", func(t *testing.T) {
		it := item{Text: "head", Line: 3}
		rest := "any of `coupons` as `c` satisfies `c`.`x` equals 1"
		expr, err := parseInlineRuleHead(rest, it, "test.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if _, ok := expr.(ast.Exists); !ok {
			t.Fatalf("expected Exists, got %T", expr)
		}
	})

	t.Run("CallHolds", func(t *testing.T) {
		it := item{Text: "head", Line: 4}
		expr, err := parseInlineRuleHead("`check`.`isSpam` holds", it, "test.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		c, ok := expr.(ast.Call)
		if !ok || c.Fn.Alias != "check" || c.Fn.Name != "isSpam" {
			t.Fatalf("expected Call check.isSpam, got %+v", expr)
		}
	})

	t.Run("CallTrailingTextError", func(t *testing.T) {
		it := item{Text: "head", Line: 5}
		_, err := parseInlineRuleHead("`check`.`isSpam` holds extra", it, "test.md")
		if err == nil || !strings.Contains(err.Error(), "unexpected trailing text after 'holds'") {
			t.Fatalf("expected trailing-text error, got %v", err)
		}
	})

	t.Run("FallThroughCompare", func(t *testing.T) {
		it := item{Text: "head", Line: 6}
		expr, err := parseInlineRuleHead("`amount` equals 5", it, "test.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		cmp, ok := expr.(ast.Compare)
		if !ok || cmp.Field != "amount" || cmp.Value != "5" {
			t.Fatalf("expected Compare amount equals 5, got %+v", expr)
		}
	})
}
