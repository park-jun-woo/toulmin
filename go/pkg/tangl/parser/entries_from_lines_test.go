//ff:func feature=tangl type=parser control=sequence
//ff:what TestEntriesFromLines — tests entriesFromLines for blank-line-skip and non-blank-line-append branches
package parser

import "testing"

func TestEntriesFromLines(t *testing.T) {
	t.Run("SkipsBlankLines", func(t *testing.T) {
		lines := []string{"- a", "   ", "- b"}
		entries := entriesFromLines(lines, 10)
		if len(entries) != 2 {
			t.Fatalf("expected 2 entries (blank line skipped), got %+v", entries)
		}
		if entries[0].Text != "- a" || entries[0].Line != 10 {
			t.Fatalf("expected first entry {- a, line 10}, got %+v", entries[0])
		}
		if entries[1].Text != "- b" || entries[1].Line != 12 {
			t.Fatalf("expected second entry {- b, line 12}, got %+v", entries[1])
		}
	})

	t.Run("AllBlank", func(t *testing.T) {
		lines := []string{"", "   ", "\t"}
		entries := entriesFromLines(lines, 1)
		if len(entries) != 0 {
			t.Fatalf("expected no entries for all-blank lines, got %+v", entries)
		}
	})

	t.Run("IndentTracked", func(t *testing.T) {
		lines := []string{"  - indented"}
		entries := entriesFromLines(lines, 5)
		if len(entries) != 1 {
			t.Fatalf("expected 1 entry, got %+v", entries)
		}
		if entries[0].Indent != 2 {
			t.Fatalf("expected indent 2, got %d", entries[0].Indent)
		}
		if entries[0].Text != "- indented" {
			t.Fatalf("expected trimmed text '- indented', got %q", entries[0].Text)
		}
	})
}
