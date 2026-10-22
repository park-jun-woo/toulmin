//ff:func feature=tangl type=validator control=sequence
//ff:what TestCheckExistsBinders — tests checkExistsBinders for valid, empty-binder, and shadowing branches
package validate

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestCheckExistsBinders(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		doc := &ast.Document{Rules: []ast.InlineRule{{Name: "r", Line: 1, Cond: ast.Exists{
			Source: ast.Path{Steps: []string{"coupons"}},
			Binder: "c",
			Cond:   ast.Compare{Field: "c", Op: "is empty"},
		}}}}
		if errs := checkExistsBinders(doc); len(errs) != 0 {
			t.Fatalf("expected no errors, got %v", errs)
		}
	})

	t.Run("NestedDistinctBindersOk", func(t *testing.T) {
		doc := &ast.Document{Rules: []ast.InlineRule{{Name: "r", Line: 1, Cond: ast.Exists{
			Source: ast.Path{Steps: []string{"orders"}},
			Binder: "o",
			Cond: ast.Exists{
				Source: ast.Path{Steps: []string{"o", "items"}},
				Binder: "i",
				Cond:   ast.Compare{Field: "i", Op: "is empty"},
			},
		}}}}
		if errs := checkExistsBinders(doc); len(errs) != 0 {
			t.Fatalf("expected no errors for distinct nested binders, got %v", errs)
		}
	})

	t.Run("EmptyBinder", func(t *testing.T) {
		doc := &ast.Document{Rules: []ast.InlineRule{{Name: "r", Line: 2, Cond: ast.Exists{
			Source: ast.Path{Steps: []string{"coupons"}},
			Binder: "",
			Cond:   ast.Compare{Field: "x", Op: "is empty"},
		}}}}
		errs := checkExistsBinders(doc)
		if len(errs) != 1 || !strings.Contains(errs[0].Error(), "empty binder name") {
			t.Fatalf("expected empty binder error, got %v", errs)
		}
	})

	t.Run("Shadowing", func(t *testing.T) {
		doc := &ast.Document{Rules: []ast.InlineRule{{Name: "r", Line: 3, Cond: ast.Exists{
			Source: ast.Path{Steps: []string{"orders"}},
			Binder: "c",
			Cond: ast.Exists{
				Source: ast.Path{Steps: []string{"c", "items"}},
				Binder: "c",
				Cond:   ast.Compare{Field: "c", Op: "is empty"},
			},
		}}}}
		errs := checkExistsBinders(doc)
		if len(errs) != 1 || !strings.Contains(errs[0].Error(), "shadows an enclosing binder") {
			t.Fatalf("expected shadowing error, got %v", errs)
		}
	})

	t.Run("UnderLogic", func(t *testing.T) {
		doc := &ast.Document{Rules: []ast.InlineRule{{Name: "r", Line: 4, Cond: ast.Logic{
			Op: "and",
			Terms: []ast.Expr{
				ast.Compare{Field: "x", Op: "is empty"},
				ast.Exists{Source: ast.Path{Steps: []string{"coupons"}}, Binder: "", Cond: ast.Compare{Field: "y", Op: "is empty"}},
			},
		}}}}
		errs := checkExistsBinders(doc)
		if len(errs) != 1 || !strings.Contains(errs[0].Error(), "empty binder name") {
			t.Fatalf("expected empty binder error under Logic, got %v", errs)
		}
	})

	t.Run("UnderEither", func(t *testing.T) {
		doc := &ast.Document{Rules: []ast.InlineRule{{Name: "r", Line: 5, Cond: ast.Either{
			Terms: []ast.Expr{
				ast.Exists{Source: ast.Path{Steps: []string{"coupons"}}, Binder: "", Cond: ast.Compare{Field: "y", Op: "is empty"}},
			},
		}}}}
		errs := checkExistsBinders(doc)
		if len(errs) != 1 || !strings.Contains(errs[0].Error(), "empty binder name") {
			t.Fatalf("expected empty binder error under Either, got %v", errs)
		}
	})

	t.Run("UnderNot", func(t *testing.T) {
		doc := &ast.Document{Rules: []ast.InlineRule{{Name: "r", Line: 6, Cond: ast.Not{
			Term: ast.Exists{Source: ast.Path{Steps: []string{"coupons"}}, Binder: "", Cond: ast.Compare{Field: "y", Op: "is empty"}},
		}}}}
		errs := checkExistsBinders(doc)
		if len(errs) != 1 || !strings.Contains(errs[0].Error(), "empty binder name") {
			t.Fatalf("expected empty binder error under Not, got %v", errs)
		}
	})
}
