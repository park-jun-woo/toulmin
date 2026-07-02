//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseRoleExpr — tests parseRoleExpr for general/counter/except/unmatched branches
package parser

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestParseRoleExpr(t *testing.T) {
	t.Run("General", func(t *testing.T) {
		role, rest, ok := parseRoleExpr("general rule extra")
		if !ok {
			t.Fatal("expected ok=true")
		}
		if role != ast.GeneralRule {
			t.Errorf("expected GeneralRule, got %v", role)
		}
		if rest != " extra" {
			t.Errorf("expected rest=' extra', got %q", rest)
		}
	})

	t.Run("Counter", func(t *testing.T) {
		role, rest, ok := parseRoleExpr("counter rule extra")
		if !ok {
			t.Fatal("expected ok=true")
		}
		if role != ast.CounterRule {
			t.Errorf("expected CounterRule, got %v", role)
		}
		if rest != " extra" {
			t.Errorf("expected rest=' extra', got %q", rest)
		}
	})

	t.Run("Except", func(t *testing.T) {
		role, rest, ok := parseRoleExpr("except rule extra")
		if !ok {
			t.Fatal("expected ok=true")
		}
		if role != ast.ExceptRule {
			t.Errorf("expected ExceptRule, got %v", role)
		}
		if rest != " extra" {
			t.Errorf("expected rest=' extra', got %q", rest)
		}
	})

	t.Run("Unmatched", func(t *testing.T) {
		role, rest, ok := parseRoleExpr("strange rule")
		if ok {
			t.Fatal("expected ok=false for an unrecognized role")
		}
		if role != 0 {
			t.Errorf("expected role=0, got %v", role)
		}
		if rest != "strange rule" {
			t.Errorf("expected rest='strange rule', got %q", rest)
		}
	})
}
