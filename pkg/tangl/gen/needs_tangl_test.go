//ff:func feature=tangl type=codegen control=iteration dimension=2
//ff:what TestNeedsTangl — tests needsTangl for provides-requires, provides-runs, internals, once-exec, undo-exec, and no-match branches
package gen

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestNeedsTangl(t *testing.T) {
	tests := []struct {
		name string
		doc  *ast.Document
		want bool
	}{
		{
			name: "empty document",
			doc:  &ast.Document{},
			want: false,
		},
		{
			name: "provide with requires triggers true",
			doc: &ast.Document{
				Provides: []ast.Endpoint{
					{Name: "ep1", Requires: []ast.Require{{}}},
				},
			},
			want: true,
		},
		{
			name: "provide with runs triggers true",
			doc: &ast.Document{
				Provides: []ast.Endpoint{
					{Name: "ep1", Runs: []string{"caseA"}},
				},
			},
			want: true,
		},
		{
			name: "provide without requires or runs is skipped",
			doc: &ast.Document{
				Provides: []ast.Endpoint{
					{Name: "ep1"},
				},
			},
			want: false,
		},
		{
			name: "internals present triggers true",
			doc: &ast.Document{
				Provides: []ast.Endpoint{{Name: "ep1"}},
				Internals: []ast.Internal{
					{Event: "start"},
				},
			},
			want: true,
		},
		{
			name: "once exec triggers true",
			doc: &ast.Document{
				Cases: []ast.Case{
					{Execs: []ast.Exec{{Kind: ast.DoExec, Once: true}}},
				},
			},
			want: true,
		},
		{
			name: "undo exec triggers true",
			doc: &ast.Document{
				Cases: []ast.Case{
					{Execs: []ast.Exec{{Kind: ast.UndoExec}}},
				},
			},
			want: true,
		},
		{
			name: "do exec without once is not enough",
			doc: &ast.Document{
				Cases: []ast.Case{
					{Execs: []ast.Exec{{Kind: ast.DoExec}}},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := needsTangl(tt.doc)
			if got != tt.want {
				t.Errorf("needsTangl() = %v, want %v", got, tt.want)
			}
		})
	}
}
