//ff:func feature=tangl type=codegen control=sequence
//ff:what TestBuildAliasImports — tests buildAliasImports for no-alias, seen-alias, and unseen-alias branches
package gen

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestBuildAliasImports(t *testing.T) {
	t.Run("no aliases", func(t *testing.T) {
		doc := &ast.Document{}
		specs, err := buildAliasImports(doc)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(specs) != 0 {
			t.Errorf("expected no import specs, got %+v", specs)
		}
	})

	t.Run("seen and unseen aliases", func(t *testing.T) {
		ref := &ast.Ref{Alias: "known", Name: "f"}
		unseenRef := &ast.Ref{Alias: "unseen", Name: "g"}
		doc := &ast.Document{
			Sees: []ast.See{
				{Alias: "known", Package: "external/pkg"},
			},
			Cases: []ast.Case{
				{
					Execs: []ast.Exec{
						{Func: ref},
						{Func: unseenRef},
					},
				},
			},
		}
		specs, err := buildAliasImports(doc)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(specs) != 2 {
			t.Fatalf("expected 2 import specs, got %+v", specs)
		}
		// aliases sorted: "known" < "unseen"
		if specs[0].Alias != "known" || specs[0].Path != "external/pkg" {
			t.Errorf("expected known alias with See path, got %+v", specs[0])
		}
		if specs[1].Alias != "unseen" || specs[1].Path != "github.com/park-jun-woo/toulmin/pkg/tangl/unseen" {
			t.Errorf("expected unseen alias with default tangl/ path, got %+v", specs[1])
		}
	})
}
