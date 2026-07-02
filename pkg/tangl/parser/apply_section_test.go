//ff:func feature=tangl type=parser control=sequence
//ff:what TestApplySection — tests applySection dispatch for the Subject, See, Definitions, and Rules section success/error branches
package parser

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestApplySection(t *testing.T) {
	t.Run("SubjectSuccess", func(t *testing.T) {
		doc := &ast.Document{}
		sec := section{Name: "Subject", HeaderLine: 1, LineOffset: 2, Lines: []string{"- this document is `myapp`"}}
		if err := applySection(doc, sec, "test.md"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if doc.Subject != "myapp" {
			t.Fatalf("expected Subject to be myapp, got %q", doc.Subject)
		}
	})

	t.Run("SubjectError", func(t *testing.T) {
		doc := &ast.Document{}
		sec := section{Name: "Subject", HeaderLine: 1, LineOffset: 2, Lines: nil}
		err := applySection(doc, sec, "test.md")
		if err == nil {
			t.Fatal("expected an error for an empty Subject section")
		}
		if !strings.Contains(err.Error(), "tangl:Subject requires") {
			t.Errorf("expected Subject requires error, got %v", err)
		}
	})

	t.Run("SeeSuccess", func(t *testing.T) {
		doc := &ast.Document{}
		sec := section{Name: "See", HeaderLine: 1, LineOffset: 2, Lines: []string{"- see `alias` from `pkg`"}}
		if err := applySection(doc, sec, "test.md"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(doc.Sees) != 1 || doc.Sees[0].Alias != "alias" {
			t.Fatalf("expected Sees to contain alias, got %+v", doc.Sees)
		}
	})

	t.Run("SeeError", func(t *testing.T) {
		doc := &ast.Document{}
		sec := section{Name: "See", HeaderLine: 1, LineOffset: 2, Lines: []string{"- bogus"}}
		err := applySection(doc, sec, "test.md")
		if err == nil {
			t.Fatal("expected an error for a malformed See item")
		}
	})

	t.Run("DefinitionsSuccess", func(t *testing.T) {
		doc := &ast.Document{}
		sec := section{Name: "Definitions", HeaderLine: 1, LineOffset: 2, Lines: []string{"- `threshold` means 5"}}
		if err := applySection(doc, sec, "test.md"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(doc.Defs) != 1 || doc.Defs[0].Name != "threshold" {
			t.Fatalf("expected Defs to contain threshold, got %+v", doc.Defs)
		}
	})

	t.Run("DefinitionsError", func(t *testing.T) {
		doc := &ast.Document{}
		sec := section{Name: "Definitions", HeaderLine: 1, LineOffset: 2, Lines: []string{"- bogus"}}
		err := applySection(doc, sec, "test.md")
		if err == nil {
			t.Fatal("expected an error for a malformed Definitions item")
		}
	})

	t.Run("RulesSuccess", func(t *testing.T) {
		doc := &ast.Document{}
		sec := section{Name: "Rules", HeaderLine: 1, LineOffset: 2, Lines: []string{"- `r1` when `amount` is empty"}}
		if err := applySection(doc, sec, "test.md"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(doc.Rules) != 1 || doc.Rules[0].Name != "r1" {
			t.Fatalf("expected Rules to contain r1, got %+v", doc.Rules)
		}
	})

	t.Run("RulesError", func(t *testing.T) {
		doc := &ast.Document{}
		sec := section{Name: "Rules", HeaderLine: 1, LineOffset: 2, Lines: []string{"- bogus"}}
		err := applySection(doc, sec, "test.md")
		if err == nil {
			t.Fatal("expected an error for a malformed Rules item")
		}
	})
}
