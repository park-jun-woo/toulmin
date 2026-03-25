//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseAttack — test attack declaration parsing
package parser

import "testing"

// TestParseAttack tests attack declaration parsing.
func TestParseAttack(t *testing.T) {
	ad, err := parseAttack("blocked attacks auth", 10, "access")
	if err != nil {
		t.Fatalf("parseAttack failed: %v", err)
	}
	if ad.Attacker != "blocked" {
		t.Errorf("expected attacker 'blocked', got %q", ad.Attacker)
	}
	if ad.Target != "auth" {
		t.Errorf("expected target 'auth', got %q", ad.Target)
	}
	if ad.Graph != "access" {
		t.Errorf("expected graph 'access', got %q", ad.Graph)
	}
	if ad.Line != 10 {
		t.Errorf("expected line 10, got %d", ad.Line)
	}
}
