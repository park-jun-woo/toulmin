//ff:func feature=tangl type=validator control=sequence
//ff:what TestRuleCallAliasErrs — tests ruleCallAliasErrs for declared, undeclared, bare, compare-operand, and exists-source calls
package validate

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestRuleCallAliasErrs(t *testing.T) {
	doc := &ast.Document{Path: "test.md"}
	aliases := map[string]bool{"check": true}

	t.Run("DeclaredAlias", func(t *testing.T) {
		r := ast.InlineRule{Name: "r", Line: 1, Cond: ast.Call{Fn: ast.Ref{Alias: "check", Name: "isSpam"}}}
		if errs := ruleCallAliasErrs(doc, r, aliases); len(errs) != 0 {
			t.Fatalf("expected no errors, got %v", errs)
		}
	})

	t.Run("BareCall", func(t *testing.T) {
		r := ast.InlineRule{Name: "r", Line: 2, Cond: ast.Call{Fn: ast.Ref{Name: "local"}}}
		if errs := ruleCallAliasErrs(doc, r, aliases); len(errs) != 0 {
			t.Fatalf("expected no errors for bare call, got %v", errs)
		}
	})

	t.Run("UndeclaredPredicate", func(t *testing.T) {
		r := ast.InlineRule{Name: "r", Line: 3, Cond: ast.Call{Fn: ast.Ref{Alias: "nlp", Name: "score"}}}
		errs := ruleCallAliasErrs(doc, r, aliases)
		if len(errs) != 1 || !strings.Contains(errs[0].Error(), "undeclared package alias") {
			t.Fatalf("expected undeclared alias error, got %v", errs)
		}
	})

	t.Run("CompareOperands", func(t *testing.T) {
		r := ast.InlineRule{Name: "r", Line: 4, Cond: ast.Compare{
			Op:  "equals",
			Lhs: ast.Call{Fn: ast.Ref{Alias: "a", Name: "f"}},
			Rhs: ast.Call{Fn: ast.Ref{Alias: "b", Name: "g"}},
		}}
		errs := ruleCallAliasErrs(doc, r, aliases)
		if len(errs) != 2 {
			t.Fatalf("expected 2 errors for Lhs and Rhs calls, got %v", errs)
		}
	})

	t.Run("CompareNonCallOperand", func(t *testing.T) {
		r := ast.InlineRule{Name: "r", Line: 5, Cond: ast.Compare{
			Op:    "equals",
			Lhs:   ast.Path{Steps: []string{"amount"}},
			Value: "5",
		}}
		if errs := ruleCallAliasErrs(doc, r, aliases); len(errs) != 0 {
			t.Fatalf("expected no errors for path operand, got %v", errs)
		}
	})

	t.Run("ExistsSource", func(t *testing.T) {
		r := ast.InlineRule{Name: "r", Line: 6, Cond: ast.Exists{
			Source: ast.Call{Fn: ast.Ref{Alias: "org", Name: "members"}},
			Binder: "m",
			Cond:   ast.Compare{Field: "m", Op: "is empty"},
		}}
		errs := ruleCallAliasErrs(doc, r, aliases)
		if len(errs) != 1 || !strings.Contains(errs[0].Error(), "org") {
			t.Fatalf("expected undeclared alias error for exists source, got %v", errs)
		}
	})

	t.Run("OtherExprKinds", func(t *testing.T) {
		r := ast.InlineRule{Name: "r", Line: 7, Cond: ast.Not{
			Term: ast.Logic{Op: "and", Terms: []ast.Expr{
				ast.Call{Fn: ast.Ref{Alias: "zzz", Name: "f"}},
			}},
		}}
		errs := ruleCallAliasErrs(doc, r, aliases)
		if len(errs) != 1 || !strings.Contains(errs[0].Error(), "zzz") {
			t.Fatalf("expected error via Not/Logic traversal, got %v", errs)
		}
	})
}
