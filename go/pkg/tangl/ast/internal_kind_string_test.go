//ff:func feature=tangl type=model control=sequence
//ff:what TestInternalKind_String — tests InternalKind.String for all known kinds and unknown default
package ast

import "testing"

func TestInternalKind_String(t *testing.T) {
	if got := OnEvent.String(); got != "OnEvent" {
		t.Errorf("OnEvent: got %q, want %q", got, "OnEvent")
	}
	if got := EveryTick.String(); got != "EveryTick" {
		t.Errorf("EveryTick: got %q, want %q", got, "EveryTick")
	}
	if got := InternalKind(99).String(); got != "InternalKind(?)" {
		t.Errorf("unknown: got %q, want %q", got, "InternalKind(?)")
	}
}
