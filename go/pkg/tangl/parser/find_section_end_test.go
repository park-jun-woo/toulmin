//ff:func feature=tangl type=parser control=sequence
//ff:what TestFindSectionEnd — tests findSectionEnd for reaching a heading line, reaching end of lines, and starting past the end
package parser

import "testing"

func TestFindSectionEnd(t *testing.T) {
	t.Run("StopsAtHeading", func(t *testing.T) {
		lines := []string{"- a", "- b", "## tangl:Next", "- c"}
		got := findSectionEnd(lines, 0)
		if got != 2 {
			t.Fatalf("expected end index 2 (heading line), got %d", got)
		}
	})

	t.Run("ReachesEndOfLines", func(t *testing.T) {
		lines := []string{"- a", "- b"}
		got := findSectionEnd(lines, 0)
		if got != len(lines) {
			t.Fatalf("expected end index %d (no heading found), got %d", len(lines), got)
		}
	})

	t.Run("StartPastEnd", func(t *testing.T) {
		lines := []string{"- a"}
		got := findSectionEnd(lines, 1)
		if got != 1 {
			t.Fatalf("expected end index 1 (start already at end), got %d", got)
		}
	})
}
