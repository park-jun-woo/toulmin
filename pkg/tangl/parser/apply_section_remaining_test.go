//ff:func feature=tangl type=parser control=sequence
//ff:what TestApplySection_Remaining — tests applySection dispatch for the Cases, Provides, and Internal section success/error branches, plus the unknown-name fallback
package parser

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestApplySection_Remaining(t *testing.T) {
	t.Run("CasesSuccess", func(t *testing.T) {
		doc := &ast.Document{}
		sec := section{Name: "Cases", HeaderLine: 1, LineOffset: 2, Lines: []string{"- in case of `c1`"}}
		if err := applySection(doc, sec, "test.md"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(doc.Cases) != 1 || doc.Cases[0].Name != "c1" {
			t.Fatalf("expected Cases to contain c1, got %+v", doc.Cases)
		}
	})

	t.Run("CasesError", func(t *testing.T) {
		doc := &ast.Document{}
		sec := section{Name: "Cases", HeaderLine: 1, LineOffset: 2, Lines: []string{"- bogus"}}
		err := applySection(doc, sec, "test.md")
		if err == nil {
			t.Fatal("expected an error for a malformed Cases item")
		}
	})

	t.Run("ProvidesSuccess", func(t *testing.T) {
		doc := &ast.Document{}
		sec := section{Name: "Provides", HeaderLine: 1, LineOffset: 2, Lines: []string{"- provides `ep1`"}}
		if err := applySection(doc, sec, "test.md"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(doc.Provides) != 1 || doc.Provides[0].Name != "ep1" {
			t.Fatalf("expected Provides to contain ep1, got %+v", doc.Provides)
		}
	})

	t.Run("ProvidesError", func(t *testing.T) {
		doc := &ast.Document{}
		sec := section{Name: "Provides", HeaderLine: 1, LineOffset: 2, Lines: []string{"- bogus"}}
		err := applySection(doc, sec, "test.md")
		if err == nil {
			t.Fatal("expected an error for a malformed Provides item")
		}
	})

	t.Run("InternalSuccess", func(t *testing.T) {
		doc := &ast.Document{}
		sec := section{Name: "Internal", HeaderLine: 1, LineOffset: 2, Lines: []string{"- on someEvent"}}
		if err := applySection(doc, sec, "test.md"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(doc.Internals) != 1 || doc.Internals[0].Event != "someEvent" {
			t.Fatalf("expected Internals to contain someEvent, got %+v", doc.Internals)
		}
	})

	t.Run("InternalError", func(t *testing.T) {
		doc := &ast.Document{}
		sec := section{Name: "Internal", HeaderLine: 1, LineOffset: 2, Lines: []string{"- bogus"}}
		err := applySection(doc, sec, "test.md")
		if err == nil {
			t.Fatal("expected an error for a malformed Internal item")
		}
	})

	t.Run("UnknownName", func(t *testing.T) {
		doc := &ast.Document{}
		sec := section{Name: "Unknown", HeaderLine: 1, LineOffset: 2, Lines: []string{"- something"}}
		if err := applySection(doc, sec, "test.md"); err != nil {
			t.Fatalf("expected no error for an unrecognized section name, got %v", err)
		}
		if doc.Subject != "" || len(doc.Sees) != 0 || len(doc.Defs) != 0 || len(doc.Rules) != 0 ||
			len(doc.Cases) != 0 || len(doc.Provides) != 0 || len(doc.Internals) != 0 {
			t.Fatalf("expected doc to remain unchanged, got %+v", doc)
		}
	})
}
