//ff:func feature=tangl type=parser control=sequence
//ff:what TestTanglSectionName — tests tanglSectionName for prefix-mismatch, empty-name, and valid-name branches
package parser

import "testing"

func TestTanglSectionName(t *testing.T) {
	t.Run("NoPrefix", func(t *testing.T) {
		name, ok := tanglSectionName("## other section")
		if ok {
			t.Fatal("expected no match when prefix does not match")
		}
		if name != "" {
			t.Fatalf("expected empty name, got %q", name)
		}
	})

	t.Run("EmptyName", func(t *testing.T) {
		name, ok := tanglSectionName("## tangl:   ")
		if ok {
			t.Fatal("expected no match when name is empty after trimming")
		}
		if name != "" {
			t.Fatalf("expected empty name, got %q", name)
		}
	})

	t.Run("Valid", func(t *testing.T) {
		name, ok := tanglSectionName("  ## tangl: rules  ")
		if !ok {
			t.Fatal("expected match")
		}
		if name != "rules" {
			t.Fatalf("expected name=%q, got %q", "rules", name)
		}
	})
}
