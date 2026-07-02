//ff:func feature=tangl type=parser control=sequence
//ff:what TestSectionOrderIndex — tests sectionOrderIndex for found and not-found branches
package parser

import "testing"

func TestSectionOrderIndex(t *testing.T) {
	t.Run("Found", func(t *testing.T) {
		idx, ok := sectionOrderIndex("Cases")
		if !ok {
			t.Fatal("expected ok=true for a known section name")
		}
		if idx != 4 {
			t.Errorf("expected idx=4, got %d", idx)
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		idx, ok := sectionOrderIndex("Bogus")
		if ok {
			t.Fatal("expected ok=false for an unknown section name")
		}
		if idx != -1 {
			t.Errorf("expected idx=-1, got %d", idx)
		}
	})
}
