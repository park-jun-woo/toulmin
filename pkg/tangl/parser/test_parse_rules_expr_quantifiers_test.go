//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseRulesExprQuantifiers — golden parse of exists/call quantifiers incl. demo C-group predicates and nested cond terms
package parser

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestParseRulesExprQuantifiers(t *testing.T) {
	rules := rulesExprByName(t)

	t.Run("C-group: has coupon (one-line exists)", func(t *testing.T) {
		ex, ok := rules["hasCoupon"].Cond.(ast.Exists)
		if !ok || ex.Binder != "c" {
			t.Fatalf("hasCoupon unexpected: %+v", rules["hasCoupon"].Cond)
		}
		src := ex.Source.(ast.Path)
		if src.Steps[0] != "coupons" {
			t.Fatalf("hasCoupon source unexpected: %+v", src)
		}
	})

	t.Run("C-group: active promotion (nested exists and-chain)", func(t *testing.T) {
		ex, ok := rules["activePromotion"].Cond.(ast.Exists)
		if !ok || ex.Binder != "p" {
			t.Fatalf("activePromotion unexpected: %+v", rules["activePromotion"].Cond)
		}
		if _, ok := ex.Cond.(ast.Logic); !ok {
			t.Fatalf("activePromotion body should be Logic, got %T", ex.Cond)
		}
	})

	t.Run("C-group: ip in list (exists with binder compare)", func(t *testing.T) {
		ex, ok := rules["ipInList"].Cond.(ast.Exists)
		if !ok || ex.Binder != "ip" {
			t.Fatalf("ipInList unexpected: %+v", rules["ipInList"].Cond)
		}
		c := ex.Cond.(ast.Compare)
		if c.Field != "ip" || c.Op != "equals" {
			t.Fatalf("ipInList body unexpected: %+v", c)
		}
		r := c.Rhs.(ast.Path)
		if r.Steps[0] != "clientIp" {
			t.Fatalf("ipInList Rhs unexpected: %+v", c.Rhs)
		}
	})

	t.Run("exists as a nested cond term", func(t *testing.T) {
		lg, ok := rules["combo"].Cond.(ast.Logic)
		if !ok || lg.Op != "or" {
			t.Fatalf("combo should be an or-Logic, got %+v", rules["combo"].Cond)
		}
		if _, ok := lg.Terms[0].(ast.Exists); !ok {
			t.Fatalf("combo first term should be Exists, got %T", lg.Terms[0])
		}
	})

	t.Run("holds as a nested cond term", func(t *testing.T) {
		lg, ok := rules["flagged"].Cond.(ast.Logic)
		if !ok {
			t.Fatalf("flagged should be Logic, got %T", rules["flagged"].Cond)
		}
		if _, ok := lg.Terms[0].(ast.Call); !ok {
			t.Fatalf("flagged first term should be Call, got %T", lg.Terms[0])
		}
	})
}
