//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseExistsTerm_Errors — tests parseExistsTerm reserved-quantifier and malformed-syntax error branches
package parser

import (
	"strings"
	"testing"
)

func TestParseExistsTerm_Errors(t *testing.T) {
	t.Run("ReservedAll", func(t *testing.T) {
		it := item{Text: "all of `coupons` as `c` satisfies `c`.`x` equals 1", Line: 5}
		_, err := parseExistsTerm(it.Text, it, "test.md")
		if err == nil || !strings.Contains(err.Error(), `quantifier "all" is reserved`) {
			t.Fatalf("expected reserved error for all, got %v", err)
		}
	})

	t.Run("ReservedCount", func(t *testing.T) {
		it := item{Text: "count of `coupons` as `c` satisfies `c`.`x` equals 1", Line: 6}
		_, err := parseExistsTerm(it.Text, it, "test.md")
		if err == nil || !strings.Contains(err.Error(), "reserved") {
			t.Fatalf("expected reserved error for count, got %v", err)
		}
	})

	t.Run("MissingOf", func(t *testing.T) {
		it := item{Text: "any `coupons` as `c` satisfies `c`.`x` equals 1", Line: 7}
		_, err := parseExistsTerm(it.Text, it, "test.md")
		if err == nil || !strings.Contains(err.Error(), "expected 'of' after 'any'") {
			t.Fatalf("expected missing-of error, got %v", err)
		}
	})

	t.Run("MissingAs", func(t *testing.T) {
		it := item{Text: "any of `coupons` `c` satisfies `c`.`x` equals 1", Line: 8}
		_, err := parseExistsTerm(it.Text, it, "test.md")
		if err == nil || !strings.Contains(err.Error(), "expected 'as'") {
			t.Fatalf("expected missing-as error, got %v", err)
		}
	})

	t.Run("MissingSatisfies", func(t *testing.T) {
		it := item{Text: "any of `coupons` as `c` `c`.`x` equals 1", Line: 9}
		_, err := parseExistsTerm(it.Text, it, "test.md")
		if err == nil || !strings.Contains(err.Error(), "expected 'satisfies'") {
			t.Fatalf("expected missing-satisfies error, got %v", err)
		}
	})

	t.Run("NoQuantifier", func(t *testing.T) {
		it := item{Text: "bogus of `coupons` as `c` satisfies `c`.`x` equals 1", Line: 10}
		_, err := parseExistsTerm(it.Text, it, "test.md")
		if err == nil || !strings.Contains(err.Error(), "expected a quantifier") {
			t.Fatalf("expected quantifier error, got %v", err)
		}
	})

	t.Run("OperandError", func(t *testing.T) {
		it := item{Text: "any of spec as `c` satisfies `c`.`x` equals 1", Line: 11}
		_, err := parseExistsTerm(it.Text, it, "test.md")
		if err == nil || !strings.Contains(err.Error(), "expected backtick-quoted spec term after 'spec'") {
			t.Fatalf("expected malformed-operand error, got %v", err)
		}
	})

	t.Run("MissingOperand", func(t *testing.T) {
		it := item{Text: "any of 42 as `c` satisfies `c`.`x` equals 1", Line: 12}
		_, err := parseExistsTerm(it.Text, it, "test.md")
		if err == nil || !strings.Contains(err.Error(), "expected a collection operand after 'any of'") {
			t.Fatalf("expected missing-operand error, got %v", err)
		}
	})

	t.Run("MissingBinder", func(t *testing.T) {
		it := item{Text: "any of `coupons` as c satisfies `c`.`x` equals 1", Line: 13}
		_, err := parseExistsTerm(it.Text, it, "test.md")
		if err == nil || !strings.Contains(err.Error(), "expected a backtick-quoted binder after 'as'") {
			t.Fatalf("expected missing-binder error, got %v", err)
		}
	})

	t.Run("OneLineCondError", func(t *testing.T) {
		it := item{Text: "any of `coupons` as `c` satisfies garbage", Line: 14}
		_, err := parseExistsTerm(it.Text, it, "test.md")
		if err == nil {
			t.Fatal("expected compare parse error, got nil")
		}
	})

	t.Run("NestedCondError", func(t *testing.T) {
		it := item{
			Text:     "any of `coupons` as `c` satisfies",
			Line:     15,
			Children: []item{{Text: "garbage", Line: 16}},
		}
		_, err := parseExistsTerm(it.Text, it, "test.md")
		if err == nil {
			t.Fatal("expected cond-list parse error, got nil")
		}
	})
}
