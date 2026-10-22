//ff:func feature=tangl type=model control=sequence
//ff:what TestRole_String — tests Role.String for all known roles and unknown default
package ast

import "testing"

func TestRole_String(t *testing.T) {
	if got := GeneralRule.String(); got != "GeneralRule" {
		t.Errorf("GeneralRule: got %q, want %q", got, "GeneralRule")
	}
	if got := CounterRule.String(); got != "CounterRule" {
		t.Errorf("CounterRule: got %q, want %q", got, "CounterRule")
	}
	if got := ExceptRule.String(); got != "ExceptRule" {
		t.Errorf("ExceptRule: got %q, want %q", got, "ExceptRule")
	}
	if got := Role(99).String(); got != "Role(?)" {
		t.Errorf("unknown: got %q, want %q", got, "Role(?)")
	}
}
