//ff:func feature=tangl type=model control=sequence
//ff:what TestDefKind_String — tests DefKind.String for all known kinds and unknown default
package ast

import "testing"

func TestDefKind_String(t *testing.T) {
	if got := ConstDef.String(); got != "ConstDef" {
		t.Errorf("ConstDef: got %q, want %q", got, "ConstDef")
	}
	if got := StructDef.String(); got != "StructDef" {
		t.Errorf("StructDef: got %q, want %q", got, "StructDef")
	}
	if got := DefKind(99).String(); got != "DefKind(?)" {
		t.Errorf("unknown: got %q, want %q", got, "DefKind(?)")
	}
}
