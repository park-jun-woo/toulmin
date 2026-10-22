//ff:func feature=tangl type=validator control=sequence
//ff:what TestCheckRuleCallAliases — tests checkRuleCallAliases for declared, undeclared, bare, and operand-position calls
package validate

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestCheckRuleCallAliases(t *testing.T) {
	t.Run("DeclaredAlias", func(t *testing.T) {
		doc := &ast.Document{
			Sees:  []ast.See{{Alias: "check", Package: "some/pkg"}},
			Rules: []ast.InlineRule{{Name: "r", Line: 1, Cond: ast.Call{Fn: ast.Ref{Alias: "check", Name: "isSpam"}}}},
		}
		if errs := checkRuleCallAliases(doc); len(errs) != 0 {
			t.Fatalf("expected no errors, got %v", errs)
		}
	})

	t.Run("UndeclaredAliasPredicate", func(t *testing.T) {
		doc := &ast.Document{
			Rules: []ast.InlineRule{{Name: "r", Line: 3, Cond: ast.Call{Fn: ast.Ref{Alias: "check", Name: "isSpam"}}}},
		}
		errs := checkRuleCallAliases(doc)
		if len(errs) != 1 || !strings.Contains(errs[0].Error(), "undeclared package alias") {
			t.Fatalf("expected undeclared alias error, got %v", errs)
		}
	})

	t.Run("BareCallOk", func(t *testing.T) {
		doc := &ast.Document{
			Rules: []ast.InlineRule{{Name: "r", Line: 1, Cond: ast.Call{Fn: ast.Ref{Name: "local"}}}},
		}
		if errs := checkRuleCallAliases(doc); len(errs) != 0 {
			t.Fatalf("expected no errors for bare call, got %v", errs)
		}
	})

	t.Run("UndeclaredAliasInOperand", func(t *testing.T) {
		doc := &ast.Document{
			Rules: []ast.InlineRule{{Name: "r", Line: 5, Cond: ast.Compare{
				Op:    "is at least",
				Lhs:   ast.Call{Fn: ast.Ref{Alias: "nlp", Name: "score"}},
				Value: "0.7",
			}}},
		}
		errs := checkRuleCallAliases(doc)
		if len(errs) != 1 || !strings.Contains(errs[0].Error(), "nlp") {
			t.Fatalf("expected undeclared alias error for operand call, got %v", errs)
		}
	})

	t.Run("UndeclaredAliasInExistsSource", func(t *testing.T) {
		doc := &ast.Document{
			Rules: []ast.InlineRule{{Name: "r", Line: 6, Cond: ast.Exists{
				Source: ast.Call{Fn: ast.Ref{Alias: "org", Name: "members"}},
				Binder: "m",
				Cond:   ast.Compare{Field: "m", Op: "is empty"},
			}}},
		}
		errs := checkRuleCallAliases(doc)
		if len(errs) != 1 || !strings.Contains(errs[0].Error(), "org") {
			t.Fatalf("expected undeclared alias error for exists source, got %v", errs)
		}
	})
}
