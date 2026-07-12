//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseCondTermDispatch — tests parseCondTerm dispatch to exists, reserved quantifier, call term, and compare fall-through
package parser

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestParseCondTermDispatch(t *testing.T) {
	t.Run("ExistsDispatch", func(t *testing.T) {
		it := item{Text: "any of `coupons` as `c` satisfies `c`.`x` equals 1", Line: 1}
		expr, err := parseCondTerm(it, "test.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if _, ok := expr.(ast.Exists); !ok {
			t.Fatalf("expected ast.Exists, got %T", expr)
		}
	})

	t.Run("ReservedQuantifierDispatch", func(t *testing.T) {
		it := item{Text: "sum of `coupons` as `c` satisfies `c`.`x` equals 1", Line: 1}
		_, err := parseCondTerm(it, "test.md")
		if err == nil || !strings.Contains(err.Error(), "reserved") {
			t.Fatalf("expected reserved error, got %v", err)
		}
	})

	t.Run("CallTermDispatch", func(t *testing.T) {
		it := item{Text: "`check`.`isSpam` holds", Line: 1}
		expr, err := parseCondTerm(it, "test.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if _, ok := expr.(ast.Call); !ok {
			t.Fatalf("expected ast.Call, got %T", expr)
		}
	})

	t.Run("CallFallThroughToCompare", func(t *testing.T) {
		it := item{Text: "`amount` equals 5", Line: 1}
		expr, err := parseCondTerm(it, "test.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if _, ok := expr.(ast.Compare); !ok {
			t.Fatalf("expected ast.Compare via fall-through, got %T", expr)
		}
	})
}
