//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseInlineRule — test inline rule parsing with expression
package parser

import "testing"

// TestParseInlineRule tests inline rule parsing with various expressions.
func TestParseInlineRule(t *testing.T) {
	ir, err := parseInlineRule(`rule "isAuthenticated" is`, "return that user is not nil", 3)
	if err != nil {
		t.Fatalf("parseInlineRule failed: %v", err)
	}
	if ir.Name != "isAuthenticated" {
		t.Errorf("expected name 'isAuthenticated', got %q", ir.Name)
	}
	if ir.Expr.Field != "user" {
		t.Errorf("expected field 'user', got %q", ir.Expr.Field)
	}
	if ir.Expr.Operator != "is not nil" {
		t.Errorf("expected operator 'is not nil', got %q", ir.Expr.Operator)
	}

	ir2, err := parseInlineRule(`rule "checkRole" is`, "return that role of spec equals admin", 5)
	if err != nil {
		t.Fatalf("parseInlineRule with of spec failed: %v", err)
	}
	if !ir2.Expr.OfSpec {
		t.Error("expected OfSpec=true")
	}
	if ir2.Expr.Field != "role" {
		t.Errorf("expected field 'role', got %q", ir2.Expr.Field)
	}
	if ir2.Expr.Operator != "equals" {
		t.Errorf("expected operator 'equals', got %q", ir2.Expr.Operator)
	}

	ir3, err := parseInlineRule(`rule "checkBoth" is`, "return that user is not nil and ip is not nil", 7)
	if err != nil {
		t.Fatalf("parseInlineRule with and failed: %v", err)
	}
	if ir3.Expr.And == nil {
		t.Fatal("expected And expression")
	}
	if ir3.Expr.And.Field != "ip" {
		t.Errorf("expected And field 'ip', got %q", ir3.Expr.And.Field)
	}
}
