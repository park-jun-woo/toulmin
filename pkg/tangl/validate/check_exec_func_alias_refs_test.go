//ff:func feature=tangl type=validator control=sequence
//ff:what TestCheckExecFuncAliasRefs — tests checkExecFuncAliasRefs for nil-func, empty-alias, declared-alias, and undeclared-alias branches
package validate

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestCheckExecFuncAliasRefs(t *testing.T) {
	t.Run("NilFunc", func(t *testing.T) {
		doc := &ast.Document{
			Cases: []ast.Case{
				{
					Name: "x",
					Execs: []ast.Exec{
						{Func: nil, Line: 1},
					},
				},
			},
		}
		errs := checkExecFuncAliasRefs(doc)
		if len(errs) != 0 {
			t.Fatalf("expected no errors, got %v", errs)
		}
	})

	t.Run("EmptyAlias", func(t *testing.T) {
		doc := &ast.Document{
			Cases: []ast.Case{
				{
					Name: "x",
					Execs: []ast.Exec{
						{Func: &ast.Ref{Name: "fn"}, Line: 1},
					},
				},
			},
		}
		errs := checkExecFuncAliasRefs(doc)
		if len(errs) != 0 {
			t.Fatalf("expected no errors, got %v", errs)
		}
	})

	t.Run("DeclaredAlias", func(t *testing.T) {
		doc := &ast.Document{
			Sees: []ast.See{
				{Alias: "pkg"},
			},
			Cases: []ast.Case{
				{
					Name: "x",
					Execs: []ast.Exec{
						{Func: &ast.Ref{Alias: "pkg", Name: "fn"}, Line: 1},
					},
				},
			},
		}
		errs := checkExecFuncAliasRefs(doc)
		if len(errs) != 0 {
			t.Fatalf("expected no errors, got %v", errs)
		}
	})

	t.Run("UndeclaredAlias", func(t *testing.T) {
		doc := &ast.Document{
			Cases: []ast.Case{
				{
					Name: "x",
					Execs: []ast.Exec{
						{Func: &ast.Ref{Alias: "missing", Name: "fn"}, Line: 4},
					},
				},
			},
		}
		errs := checkExecFuncAliasRefs(doc)
		if len(errs) != 1 {
			t.Fatalf("expected 1 error, got %v", errs)
		}
		if !strings.Contains(errs[0].Error(), "undeclared package alias") {
			t.Errorf("expected undeclared package alias error, got %v", errs[0])
		}
	})
}
