//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseFullExample — parse the full TANGL example and verify all AST fields
package parser

import "testing"

// TestParseFullExample parses the full TANGL example from the grammar spec.
func TestParseFullExample(t *testing.T) {
	f, err := ParseString(templateFullExample)
	if err != nil {
		t.Fatalf("ParseString failed: %v", err)
	}

	if len(f.Imports) != 1 {
		t.Fatalf("expected 1 import, got %d", len(f.Imports))
	}
	if f.Imports[0].Alias != "policy" {
		t.Errorf("expected import alias 'policy', got %q", f.Imports[0].Alias)
	}
	if f.Imports[0].Package != "github.com/park-jun-woo/toulmin/pkg/policy" {
		t.Errorf("expected import package, got %q", f.Imports[0].Package)
	}

	if len(f.Rules) != 1 {
		t.Fatalf("expected 1 inline rule, got %d", len(f.Rules))
	}
	if f.Rules[0].Name != "isAuthenticated" {
		t.Errorf("expected rule name 'isAuthenticated', got %q", f.Rules[0].Name)
	}
	if f.Rules[0].Expr.Field != "user" {
		t.Errorf("expected expr field 'user', got %q", f.Rules[0].Expr.Field)
	}
	if f.Rules[0].Expr.Operator != "is not nil" {
		t.Errorf("expected operator 'is not nil', got %q", f.Rules[0].Expr.Operator)
	}

	if len(f.Graphs) != 1 {
		t.Fatalf("expected 1 graph, got %d", len(f.Graphs))
	}
	if f.Graphs[0].Name != "access" {
		t.Errorf("expected graph name 'access', got %q", f.Graphs[0].Name)
	}
	if f.Graphs[0].ID != "api:access" {
		t.Errorf("expected graph id 'api:access', got %q", f.Graphs[0].ID)
	}

	if len(f.Bindings) != 4 {
		t.Fatalf("expected 4 bindings, got %d", len(f.Bindings))
	}

	auth := f.Bindings[0]
	if auth.Name != "auth" || auth.Role != "rule" || auth.Func != "isAuthenticated" || auth.Graph != "access" {
		t.Errorf("auth binding mismatch: %+v", auth)
	}

	admin := f.Bindings[1]
	if admin.Name != "admin" || admin.Role != "rule" || admin.Func != "policy.IsInRole" || admin.Graph != "access" {
		t.Errorf("admin binding mismatch: %+v", admin)
	}
	if len(admin.Specs) != 1 || admin.Specs[0].Name != "policy.Role" {
		t.Errorf("admin specs mismatch: %+v", admin.Specs)
	}
	if len(admin.Specs[0].Args) != 1 {
		t.Errorf("expected 1 arg for policy.Role, got %d", len(admin.Specs[0].Args))
	}
	if admin.Specs[0].Args[0] != "admin" {
		t.Errorf("expected arg 'admin', got %v", admin.Specs[0].Args[0])
	}

	blocked := f.Bindings[2]
	if blocked.Name != "blocked" || blocked.Role != "counter" || blocked.Func != "policy.IsIPBlocked" {
		t.Errorf("blocked binding mismatch: %+v", blocked)
	}
	if len(blocked.Specs) != 1 || blocked.Specs[0].Name != "policy.IPList" {
		t.Errorf("blocked specs mismatch: %+v", blocked.Specs)
	}

	exempt := f.Bindings[3]
	if exempt.Name != "exempt" || exempt.Role != "except" || exempt.Func != "policy.IsInternalIP" {
		t.Errorf("exempt binding mismatch: %+v", exempt)
	}

	if len(f.Attacks) != 2 {
		t.Fatalf("expected 2 attacks, got %d", len(f.Attacks))
	}
	if f.Attacks[0].Attacker != "blocked" || f.Attacks[0].Target != "auth" {
		t.Errorf("attack 0 mismatch: %+v", f.Attacks[0])
	}
	if f.Attacks[1].Attacker != "exempt" || f.Attacks[1].Target != "blocked" {
		t.Errorf("attack 1 mismatch: %+v", f.Attacks[1])
	}

	if len(f.Evals) != 1 {
		t.Fatalf("expected 1 eval, got %d", len(f.Evals))
	}
	if f.Evals[0].Name != "acResult" || f.Evals[0].Graph != "access" || !f.Evals[0].Trace {
		t.Errorf("eval mismatch: %+v", f.Evals[0])
	}
}
