//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseCompareOperands — tests parseCompare operand left/right sides (path, spec, map key, result-of, literal)
package parser

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestParseCompareOperands(t *testing.T) {
	t.Run("SingleStepStaysField", func(t *testing.T) {
		cmp, err := parseCompare("`credit` is at least 650", 5, "test.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cmp.Field != "credit" || cmp.Lhs != nil || cmp.Rhs != nil {
			t.Fatalf("single-step left side must demote to Field with nil Lhs/Rhs, got %+v", cmp)
		}
	})

	t.Run("PathLeftSide", func(t *testing.T) {
		cmp, err := parseCompare("`author`.`trustScore` is at least 0.7", 6, "test.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		p, ok := cmp.Lhs.(ast.Path)
		if !ok || len(p.Steps) != 2 || cmp.Field != "" || cmp.Value != "0.7" {
			t.Fatalf("expected Path Lhs, empty Field, Value=0.7, got %+v", cmp)
		}
	})

	t.Run("SpecLeftSide", func(t *testing.T) {
		cmp, err := parseCompare("spec `discount`.`min price` is at most `basePrice`", 7, "test.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if _, ok := cmp.Lhs.(ast.SpecRef); !ok {
			t.Fatalf("expected SpecRef Lhs, got %+v", cmp)
		}
		rp, ok := cmp.Rhs.(ast.Path)
		if !ok || len(rp.Steps) != 1 || rp.Steps[0] != "basePrice" || cmp.Value != "" {
			t.Fatalf("expected Path Rhs=basePrice and empty Value, got %+v", cmp)
		}
	})

	t.Run("MapKeyLeftSide", func(t *testing.T) {
		cmp, err := parseCompare("`attributes`.`internal` equals true", 8, "test.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		p, ok := cmp.Lhs.(ast.Path)
		if !ok || len(p.Steps) != 2 || cmp.Value != "true" || cmp.Rhs != nil {
			t.Fatalf("expected 2-step Path Lhs, Value=true, nil Rhs, got %+v", cmp)
		}
	})

	t.Run("ResultOfLeftSide", func(t *testing.T) {
		cmp, err := parseCompare("result of `nlp`.`score` is at least 0.7", 9, "test.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		c, ok := cmp.Lhs.(ast.Call)
		if !ok || c.Fn.Alias != "nlp" || c.Fn.Name != "score" || cmp.Value != "0.7" {
			t.Fatalf("expected Call Lhs nlp.score, Value=0.7, got %+v", cmp)
		}
	})

	t.Run("LiteralRightSideUnchanged", func(t *testing.T) {
		cmp, err := parseCompare("`grade` is in a, b, c", 10, "test.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cmp.Field != "grade" || cmp.Op != "is in" || cmp.Value != "a, b, c" || cmp.Rhs != nil {
			t.Fatalf("literal right side must stay Value string, got %+v", cmp)
		}
	})
}
