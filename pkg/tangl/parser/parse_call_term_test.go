//ff:func feature=tangl type=parser control=sequence
//ff:what TestTryParseCallTerm — tests tryParseCallTerm for holds success, fall-through, and trailing-text error branches
package parser

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestTryParseCallTerm(t *testing.T) {
	t.Run("HoldsSuccess", func(t *testing.T) {
		expr, consumed, err := tryParseCallTerm("`check`.`isSpam` holds", 1, "test.md")
		if err != nil || !consumed {
			t.Fatalf("expected consumed with no error, got consumed=%v err=%v", consumed, err)
		}
		c, ok := expr.(ast.Call)
		if !ok || c.Fn.Alias != "check" || c.Fn.Name != "isSpam" {
			t.Fatalf("expected Call check.isSpam, got %+v", expr)
		}
	})

	t.Run("FallThroughCompare", func(t *testing.T) {
		expr, consumed, err := tryParseCallTerm("`amount` equals 5", 1, "test.md")
		if consumed || err != nil {
			t.Fatalf("expected fall-through (not consumed, no error), got consumed=%v err=%v", consumed, err)
		}
		if expr != nil {
			t.Errorf("expected nil expr on fall-through, got %+v", expr)
		}
	})

	t.Run("FallThroughNoBacktick", func(t *testing.T) {
		_, consumed, err := tryParseCallTerm("bogus holds", 1, "test.md")
		if consumed || err != nil {
			t.Fatalf("expected fall-through for non-backtick text, got consumed=%v err=%v", consumed, err)
		}
	})

	t.Run("TrailingText", func(t *testing.T) {
		_, consumed, err := tryParseCallTerm("`check`.`isSpam` holds extra", 1, "test.md")
		if !consumed || err == nil || !strings.Contains(err.Error(), "unexpected trailing text after 'holds'") {
			t.Fatalf("expected trailing-text error, got consumed=%v err=%v", consumed, err)
		}
	})
}
