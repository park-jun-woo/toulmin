//ff:func feature=tangl type=codegen control=iteration dimension=1
//ff:what TestRoleMethod — tests roleMethod for counter, except, and default (general) role branches
package gen

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestRoleMethod(t *testing.T) {
	tests := []struct {
		name string
		in   ast.Role
		want string
	}{
		{"counter role", ast.CounterRule, "Counter"},
		{"except role", ast.ExceptRule, "Except"},
		{"general role defaults to Rule", ast.GeneralRule, "Rule"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := roleMethod(tt.in)
			if got != tt.want {
				t.Errorf("roleMethod(%v) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}
