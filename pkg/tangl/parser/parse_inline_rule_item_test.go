//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseInlineRuleItem — tests parseInlineRuleItem for name/when/condList/compare branches
package parser

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestParseInlineRuleItem(t *testing.T) {
	t.Run("NoBacktickName", func(t *testing.T) {
		it := item{Text: "rule1 when `x` equals 1", Line: 1}
		_, err := parseInlineRuleItem(it, "test.md")
		if err == nil {
			t.Fatal("expected an error for missing backtick-quoted rule name")
		}
		if !strings.Contains(err.Error(), "expected backtick-quoted rule name") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("NoWhen", func(t *testing.T) {
		it := item{Text: "`rule1` if `x` equals 1", Line: 2}
		_, err := parseInlineRuleItem(it, "test.md")
		if err == nil {
			t.Fatal("expected an error for missing 'when'")
		}
		if !strings.Contains(err.Error(), "expected 'when' after rule name") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("CondListError", func(t *testing.T) {
		it := item{Text: "`rule1` when", Line: 3}
		_, err := parseInlineRuleItem(it, "test.md")
		if err == nil {
			t.Fatal("expected an error from parseCondList for missing children")
		}
		if !strings.Contains(err.Error(), "expected at least one condition") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("CondListSuccess", func(t *testing.T) {
		it := item{
			Text: "`rule1` when",
			Line: 4,
			Children: []item{
				{Text: "`x` equals 1", Line: 5},
			},
		}
		rule, err := parseInlineRuleItem(it, "test.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if rule.Name != "rule1" {
			t.Errorf("expected Name=rule1, got %q", rule.Name)
		}
		cmp, ok := rule.Cond.(ast.Compare)
		if !ok {
			t.Fatalf("expected ast.Compare cond, got %T", rule.Cond)
		}
		if cmp.Field != "x" || cmp.Op != "equals" || cmp.Value != "1" {
			t.Errorf("unexpected compare: %+v", cmp)
		}
	})

	t.Run("CompareError", func(t *testing.T) {
		it := item{Text: "`rule1` when notbacktick equals 1", Line: 6}
		_, err := parseInlineRuleItem(it, "test.md")
		if err == nil {
			t.Fatal("expected an error from parseCompare for a missing backtick field")
		}
		if !strings.Contains(err.Error(), "expected backtick-quoted field") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("CompareSuccess", func(t *testing.T) {
		it := item{Text: "`rule1` when `x` equals 1", Line: 7}
		rule, err := parseInlineRuleItem(it, "test.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if rule.Name != "rule1" {
			t.Errorf("expected Name=rule1, got %q", rule.Name)
		}
		cmp, ok := rule.Cond.(ast.Compare)
		if !ok {
			t.Fatalf("expected ast.Compare cond, got %T", rule.Cond)
		}
		if cmp.Field != "x" || cmp.Op != "equals" || cmp.Value != "1" {
			t.Errorf("unexpected compare: %+v", cmp)
		}
	})
}
