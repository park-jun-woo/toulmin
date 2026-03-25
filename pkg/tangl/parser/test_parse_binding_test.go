//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseBinding — test rule binding parsing with specs and qualifier
package parser

import "testing"

// TestParseBinding tests binding parsing with various forms.
func TestParseBinding(t *testing.T) {
	rb, err := parseBinding("auth is a rule using isAuthenticated", 10, "access", 0)
	if err != nil {
		t.Fatalf("parseBinding simple failed: %v", err)
	}
	if rb.Name != "auth" || rb.Role != "rule" || rb.Func != "isAuthenticated" || rb.Graph != "access" {
		t.Errorf("simple binding mismatch: %+v", rb)
	}

	rb2, err := parseBinding("admin is a rule of access using policy.IsInRole with policy.Role(\"admin\")", 11, "", 0)
	if err != nil {
		t.Fatalf("parseBinding with graph and spec failed: %v", err)
	}
	if rb2.Graph != "access" {
		t.Errorf("expected graph 'access', got %q", rb2.Graph)
	}
	if rb2.Func != "policy.IsInRole" {
		t.Errorf("expected func 'policy.IsInRole', got %q", rb2.Func)
	}
	if len(rb2.Specs) != 1 {
		t.Fatalf("expected 1 spec, got %d", len(rb2.Specs))
	}
	if rb2.Specs[0].Name != "policy.Role" {
		t.Errorf("expected spec name 'policy.Role', got %q", rb2.Specs[0].Name)
	}

	rb3, err := parseBinding("w is a rule of g using fn qualified 0.8", 12, "", 0)
	if err != nil {
		t.Fatalf("parseBinding with qualifier failed: %v", err)
	}
	if rb3.Qualifier != 0.8 {
		t.Errorf("expected qualifier 0.8, got %f", rb3.Qualifier)
	}

	rb4, err := parseBinding("blocked is a counter using policy.IsIPBlocked with policy.IPList(\"blocklist\")", 13, "access", 0)
	if err != nil {
		t.Fatalf("parseBinding counter failed: %v", err)
	}
	if rb4.Role != "counter" {
		t.Errorf("expected role 'counter', got %q", rb4.Role)
	}

	rb5, err := parseBinding("exempt is an except using policy.IsInternalIP", 14, "access", 0)
	if err != nil {
		t.Fatalf("parseBinding except failed: %v", err)
	}
	if rb5.Role != "except" {
		t.Errorf("expected role 'except', got %q", rb5.Role)
	}
}
