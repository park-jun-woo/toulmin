//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestRoleRoundTrip — Role↔RoleCode golden values, roundtrip identity, unknown errors
package codebook

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestRoleRoundTrip(t *testing.T) {
	cases := []struct {
		role ast.Role
		code RoleCode
	}{
		{ast.GeneralRule, RoleGeneral},
		{ast.CounterRule, RoleCounter},
		{ast.ExceptRule, RoleExcept},
	}
	for _, c := range cases {
		code, err := RoleFromAST(c.role)
		if err != nil || code != c.code {
			t.Fatalf("RoleFromAST(%v): got %d err=%v, want %d", c.role, uint16(code), err, uint16(c.code))
		}
		if uint16(c.code) > 2 {
			t.Errorf("RoleCode %d exceeds the 2-bit range", uint16(c.code))
		}
		role, err := RoleToAST(c.code)
		if err != nil || role != c.role {
			t.Fatalf("RoleToAST(%d): got %v err=%v, want %v", uint16(c.code), role, err, c.role)
		}
	}
	if _, err := RoleFromAST(ast.Role(99)); err == nil {
		t.Error("RoleFromAST(unknown): expected error")
	}
	if _, err := RoleToAST(RoleCode(99)); err == nil {
		t.Error("RoleToAST(unknown): expected error")
	}
}
