//ff:func feature=tangl type=codegen control=sequence
//ff:what TestCollectUsedAliases — tests collectUsedAliases across all nil/empty/non-empty branch combinations
package gen

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestCollectUsedAliases(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		used := collectUsedAliases(&ast.Document{})
		if len(used) != 0 {
			t.Errorf("expected empty set, got %+v", used)
		}
	})

	t.Run("defs specRef nil and empty alias", func(t *testing.T) {
		doc := &ast.Document{
			Defs: []ast.Definition{
				{Name: "noSpecRef"}, // SpecRef == nil
				{Name: "emptyAlias", SpecRef: &ast.Ref{Name: "x"}}, // SpecRef != nil, Alias == ""
			},
		}
		used := collectUsedAliases(doc)
		if len(used) != 0 {
			t.Errorf("expected empty set, got %+v", used)
		}
	})

	t.Run("defs specRef with alias", func(t *testing.T) {
		doc := &ast.Document{
			Defs: []ast.Definition{
				{Name: "withAlias", SpecRef: &ast.Ref{Alias: "pkg", Name: "x"}},
			},
		}
		used := collectUsedAliases(doc)
		if !used["pkg"] || len(used) != 1 {
			t.Errorf("expected {pkg: true}, got %+v", used)
		}
	})

	t.Run("nodes using nil and empty alias", func(t *testing.T) {
		doc := &ast.Document{
			Cases: []ast.Case{
				{
					Nodes: []ast.Node{
						{Name: "n1"}, // Using == nil
						{Name: "n2", Using: &ast.Ref{Name: "localFn"}}, // Using != nil, Alias == ""
					},
				},
			},
		}
		used := collectUsedAliases(doc)
		if len(used) != 0 {
			t.Errorf("expected empty set, got %+v", used)
		}
	})

	t.Run("nodes using with alias", func(t *testing.T) {
		doc := &ast.Document{
			Cases: []ast.Case{
				{
					Nodes: []ast.Node{
						{Name: "n1", Using: &ast.Ref{Alias: "svc", Name: "check"}},
					},
				},
			},
		}
		used := collectUsedAliases(doc)
		if !used["svc"] || len(used) != 1 {
			t.Errorf("expected {svc: true}, got %+v", used)
		}
	})

	t.Run("execs func nil and empty alias", func(t *testing.T) {
		doc := &ast.Document{
			Cases: []ast.Case{
				{
					Execs: []ast.Exec{
						{Kind: ast.RunExec, Case: "other"},                  // Func == nil
						{Kind: ast.DoExec, Func: &ast.Ref{Name: "localFn"}}, // Func != nil, Alias == ""
					},
				},
			},
		}
		used := collectUsedAliases(doc)
		if len(used) != 0 {
			t.Errorf("expected empty set, got %+v", used)
		}
	})

	t.Run("execs func with alias", func(t *testing.T) {
		doc := &ast.Document{
			Cases: []ast.Case{
				{
					Execs: []ast.Exec{
						{Kind: ast.DoExec, Func: &ast.Ref{Alias: "bank", Name: "withdraw"}},
					},
				},
			},
		}
		used := collectUsedAliases(doc)
		if !used["bank"] || len(used) != 1 {
			t.Errorf("expected {bank: true}, got %+v", used)
		}
	})
}
