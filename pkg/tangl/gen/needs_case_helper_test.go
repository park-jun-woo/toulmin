//ff:func feature=tangl type=codegen control=iteration dimension=2
//ff:what TestNeedsCaseHelper — tests needsCaseHelper for checking-node, until-internal, and no-match branches
package gen

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestNeedsCaseHelper(t *testing.T) {
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
			name: "node without checking and internal without until",
			doc: &ast.Document{
				Cases: []ast.Case{
					{Nodes: []ast.Node{{Name: "a"}}},
				},
				Internals: []ast.Internal{
					{Event: "start"},
				},
			},
			want: false,
		},
		{
			name: "node with checking triggers true",
			doc: &ast.Document{
				Cases: []ast.Case{
					{Nodes: []ast.Node{{Name: "a"}, {Name: "b", Checking: "otherCase"}}},
				},
			},
			want: true,
		},
		{
			name: "internal with until triggers true",
			doc: &ast.Document{
				Cases: []ast.Case{
					{Nodes: []ast.Node{{Name: "a"}}},
				},
				Internals: []ast.Internal{
					{Event: "start", Until: "doneCase"},
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := needsCaseHelper(tt.doc)
			if got != tt.want {
				t.Errorf("needsCaseHelper() = %v, want %v", got, tt.want)
			}
		})
	}
}
