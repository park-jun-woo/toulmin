//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseOperand — tests parseOperand for path, spec, result-of-call, non-operand, and malformed branches
package parser

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestParseOperand(t *testing.T) {
	t.Run("Path", func(t *testing.T) {
		op, _, ok, err := parseOperand("`author`.`trustScore` is at least 0.7", 1, "test.md")
		if err != nil || !ok {
			t.Fatalf("unexpected ok=%v err=%v", ok, err)
		}
		p, isPath := op.(ast.Path)
		if !isPath || len(p.Steps) != 2 {
			t.Fatalf("expected Path with 2 steps, got %+v", op)
		}
	})

	t.Run("SpecRef", func(t *testing.T) {
		op, _, ok, err := parseOperand("spec `discount`.`min price` is at most `basePrice`", 1, "test.md")
		if err != nil || !ok {
			t.Fatalf("unexpected ok=%v err=%v", ok, err)
		}
		s, isSpec := op.(ast.SpecRef)
		if !isSpec || len(s.Steps) != 2 || s.Steps[0] != "discount" || s.Steps[1] != "min price" {
			t.Fatalf("expected SpecRef [discount min price], got %+v", op)
		}
	})

	t.Run("ResultOfCall", func(t *testing.T) {
		op, _, ok, err := parseOperand("result of `nlp`.`score` is at least 0.7", 1, "test.md")
		if err != nil || !ok {
			t.Fatalf("unexpected ok=%v err=%v", ok, err)
		}
		c, isCall := op.(ast.Call)
		if !isCall || c.Fn.Alias != "nlp" || c.Fn.Name != "score" {
			t.Fatalf("expected Call nlp.score, got %+v", op)
		}
	})

	t.Run("NotAnOperand", func(t *testing.T) {
		op, rest, ok, err := parseOperand("650", 1, "test.md")
		if ok || err != nil {
			t.Fatalf("expected ok=false err=nil, got ok=%v err=%v", ok, err)
		}
		if op != nil || rest != "650" {
			t.Errorf("expected nil operand and original rest, got %v / %q", op, rest)
		}
	})

	t.Run("SpecMissingPath", func(t *testing.T) {
		_, _, _, err := parseOperand("spec bogus", 1, "test.md")
		if err == nil || !strings.Contains(err.Error(), "expected backtick-quoted spec term") {
			t.Fatalf("expected spec-term error, got %v", err)
		}
	})

	t.Run("ResultMissingOf", func(t *testing.T) {
		_, _, _, err := parseOperand("result `nlp`.`score`", 1, "test.md")
		if err == nil || !strings.Contains(err.Error(), "expected 'of' after 'result'") {
			t.Fatalf("expected missing-of error, got %v", err)
		}
	})

	t.Run("ResultMissingRef", func(t *testing.T) {
		_, _, _, err := parseOperand("result of bogus", 1, "test.md")
		if err == nil || !strings.Contains(err.Error(), "expected backtick-quoted function reference") {
			t.Fatalf("expected missing-ref error, got %v", err)
		}
	})
}
