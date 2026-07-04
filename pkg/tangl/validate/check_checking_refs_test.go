//ff:func feature=tangl type=validator control=sequence
//ff:what TestCheckCheckingRefs — tests checkCheckingRefs for empty-checking-skip, valid-target, and missing-target branches via subtests
package validate

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestCheckCheckingRefs(t *testing.T) {
	t.Run("NoCases", func(t *testing.T) {
		doc := &ast.Document{}
		errs := checkCheckingRefs(doc)
		if len(errs) != 0 {
			t.Fatalf("expected no errors, got %v", errs)
		}
	})

	t.Run("SkipsEmptyChecking", func(t *testing.T) {
		doc := &ast.Document{
			Cases: []ast.Case{
				{
					Name: "x",
					Nodes: []ast.Node{
						{Name: "a", Checking: nil},
					},
				},
			},
		}
		errs := checkCheckingRefs(doc)
		if len(errs) != 0 {
			t.Fatalf("expected no errors, got %v", errs)
		}
	})

	t.Run("ValidTarget", func(t *testing.T) {
		doc := &ast.Document{
			Cases: []ast.Case{
				{
					Name: "x",
					Nodes: []ast.Node{
						{Name: "a", Checking: &ast.Ref{Name: "y"}},
					},
				},
				{Name: "y"},
			},
		}
		errs := checkCheckingRefs(doc)
		if len(errs) != 0 {
			t.Fatalf("expected no errors, got %v", errs)
		}
	})

	t.Run("SkipsAliasQualified", func(t *testing.T) {
		// An alias-qualified checking target lives in another document; its
		// existence check is deferred to Phase015's loader, so checkCheckingRefs
		// must NOT flag it even though no local case named "elsewhere" exists.
		doc := &ast.Document{
			Cases: []ast.Case{
				{
					Name: "x",
					Nodes: []ast.Node{
						{Name: "a", Checking: &ast.Ref{Alias: "gate", Name: "elsewhere"}, Line: 3},
					},
				},
			},
		}
		errs := checkCheckingRefs(doc)
		if len(errs) != 0 {
			t.Fatalf("expected no errors for alias-qualified checking, got %v", errs)
		}
	})

	t.Run("MissingTarget", func(t *testing.T) {
		doc := &ast.Document{
			Cases: []ast.Case{
				{
					Name: "x",
					Nodes: []ast.Node{
						{Name: "a", Checking: &ast.Ref{Name: "missing"}, Line: 5},
					},
				},
			},
		}
		errs := checkCheckingRefs(doc)
		if len(errs) != 1 {
			t.Fatalf("expected 1 error, got %v", errs)
		}
		if !strings.Contains(errs[0].Error(), "does not exist") {
			t.Errorf("expected does-not-exist error, got %v", errs[0])
		}
	})
}
