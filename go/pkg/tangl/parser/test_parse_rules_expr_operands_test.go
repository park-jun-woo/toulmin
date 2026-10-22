//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseRulesExprOperands — golden parse of operand forms (plain, field path, map key, spec, result-of, holds)
package parser

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestParseRulesExprOperands(t *testing.T) {
	rules := rulesExprByName(t)

	t.Run("plain compare stays Field", func(t *testing.T) {
		c := rules["credit ok"].Cond.(ast.Compare)
		if c.Field != "credit score" || c.Lhs != nil || c.Value != "650" {
			t.Fatalf("credit ok unexpected: %+v", c)
		}
	})

	t.Run("field path Lhs", func(t *testing.T) {
		c := rules["trusted author"].Cond.(ast.Compare)
		p, ok := c.Lhs.(ast.Path)
		if !ok || len(p.Steps) != 2 || p.Steps[1] != "trustScore" {
			t.Fatalf("trusted author unexpected: %+v", c)
		}
	})

	t.Run("map key path", func(t *testing.T) {
		c := rules["is internal"].Cond.(ast.Compare)
		p, ok := c.Lhs.(ast.Path)
		if !ok || p.Steps[0] != "attributes" || p.Steps[1] != "internal" || c.Value != "true" {
			t.Fatalf("is internal unexpected: %+v", c)
		}
	})

	t.Run("spec ref Lhs with backtick Rhs", func(t *testing.T) {
		c := rules["within budget"].Cond.(ast.Compare)
		s, ok := c.Lhs.(ast.SpecRef)
		if !ok || s.Steps[0] != "discount" || s.Steps[1] != "min price" {
			t.Fatalf("within budget Lhs unexpected: %+v", c)
		}
		r, ok := c.Rhs.(ast.Path)
		if !ok || r.Steps[0] != "basePrice" || c.Value != "" {
			t.Fatalf("within budget Rhs unexpected: %+v", c)
		}
	})

	t.Run("result of call value", func(t *testing.T) {
		c := rules["high nlp"].Cond.(ast.Compare)
		call, ok := c.Lhs.(ast.Call)
		if !ok || call.Fn.Alias != "nlp" || call.Fn.Name != "score" {
			t.Fatalf("high nlp unexpected: %+v", c)
		}
	})

	t.Run("holds predicate", func(t *testing.T) {
		call, ok := rules["spam"].Cond.(ast.Call)
		if !ok || call.Fn.Alias != "check" || call.Fn.Name != "isSpam" {
			t.Fatalf("spam unexpected: %+v", rules["spam"].Cond)
		}
	})
}
