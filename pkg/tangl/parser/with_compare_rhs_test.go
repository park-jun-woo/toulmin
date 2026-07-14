//ff:func feature=tangl type=parser control=sequence
//ff:what TestWithCompareRhs — tests withCompareRhs for literal value, operand rhs, trailing-text, and operand-error branches
package parser

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestWithCompareRhs(t *testing.T) {
	t.Run("LiteralValue", func(t *testing.T) {
		cmp, err := withCompareRhs(ast.Compare{Field: "amount", Op: "equals"}, " 650 ", 1, "test.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cmp.Value != "650" || cmp.Rhs != nil {
			t.Fatalf("expected literal Value=650 with nil Rhs, got %+v", cmp)
		}
	})

	t.Run("QuotedLiteralValue", func(t *testing.T) {
		cmp, err := withCompareRhs(ast.Compare{}, ` "VERY VIP" `, 1, "test.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cmp.Value != "VERY VIP" {
			t.Fatalf("quoted literal Value = %q, want %q", cmp.Value, "VERY VIP")
		}
	})

	t.Run("UnmatchedQuoteIsVerbatim", func(t *testing.T) {
		cmp, err := withCompareRhs(ast.Compare{}, `"VIP`, 1, "test.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cmp.Value != `"VIP` {
			t.Fatalf("unmatched literal Value = %q, want verbatim", cmp.Value)
		}
	})

	t.Run("QuotedIsInItems", func(t *testing.T) {
		cmp, err := withCompareRhs(ast.Compare{Op: "is in"}, ` "VIP", "GOLD" `, 1, "test.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cmp.Value != "VIP, GOLD" {
			t.Fatalf("is-in Value = %q, want %q", cmp.Value, "VIP, GOLD")
		}
	})

	t.Run("UnquotedIsInItems", func(t *testing.T) {
		cmp, err := withCompareRhs(ast.Compare{Op: "is in"}, `VIP, GOLD`, 1, "test.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cmp.Value != "VIP, GOLD" {
			t.Fatalf("is-in Value = %q, want backward-compatible list", cmp.Value)
		}
	})

	t.Run("OperandRhs", func(t *testing.T) {
		cmp, err := withCompareRhs(ast.Compare{Field: "min", Op: "is at most"}, "`basePrice`", 2, "test.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		p, ok := cmp.Rhs.(ast.Path)
		if !ok || len(p.Steps) != 1 || p.Steps[0] != "basePrice" {
			t.Fatalf("expected Rhs Path[basePrice], got %+v", cmp.Rhs)
		}
	})

	t.Run("TrailingText", func(t *testing.T) {
		_, err := withCompareRhs(ast.Compare{}, "`basePrice` extra", 3, "test.md")
		if err == nil || !strings.Contains(err.Error(), "unexpected trailing text after operand") {
			t.Fatalf("expected trailing-text error, got %v", err)
		}
	})

	t.Run("OperandError", func(t *testing.T) {
		_, err := withCompareRhs(ast.Compare{}, "spec bogus", 4, "test.md")
		if err == nil || !strings.Contains(err.Error(), "expected backtick-quoted spec term") {
			t.Fatalf("expected malformed-operand error, got %v", err)
		}
	})
}
