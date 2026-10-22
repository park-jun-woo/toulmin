//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseInlineRuleItemHeads — tests rule-head one-line dispatch to exists (inline/nested), holds call, and reserved quantifier
package parser

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestParseInlineRuleItemHeads(t *testing.T) {
	t.Run("ExistsHeadOneLine", func(t *testing.T) {
		it := item{Text: "`hasCoupon` when any of `coupons` as `c` satisfies `c`.`minPrice` is at most `basePrice`", Line: 8}
		rule, err := parseInlineRuleItem(it, "test.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if rule.Name != "hasCoupon" {
			t.Errorf("expected Name=hasCoupon, got %q", rule.Name)
		}
		if _, ok := rule.Cond.(ast.Exists); !ok {
			t.Fatalf("expected ast.Exists cond, got %T", rule.Cond)
		}
	})

	t.Run("ExistsHeadNestedSatisfies", func(t *testing.T) {
		it := item{
			Text: "`hasCoupon` when any of `coupons` as `c` satisfies",
			Line: 9,
			Children: []item{
				{Text: "`c`.`active` equals true", Line: 10},
			},
		}
		rule, err := parseInlineRuleItem(it, "test.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if _, ok := rule.Cond.(ast.Exists); !ok {
			t.Fatalf("expected ast.Exists cond, got %T", rule.Cond)
		}
	})

	t.Run("CallHeadHolds", func(t *testing.T) {
		it := item{Text: "`spam` when `check`.`isSpam` holds", Line: 11}
		rule, err := parseInlineRuleItem(it, "test.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if _, ok := rule.Cond.(ast.Call); !ok {
			t.Fatalf("expected ast.Call cond, got %T", rule.Cond)
		}
	})

	t.Run("ReservedQuantifierHead", func(t *testing.T) {
		it := item{Text: "`r` when all of `coupons` as `c` satisfies `c`.`x` equals 1", Line: 12}
		_, err := parseInlineRuleItem(it, "test.md")
		if err == nil || !strings.Contains(err.Error(), "reserved") {
			t.Fatalf("expected reserved error, got %v", err)
		}
	})
}
