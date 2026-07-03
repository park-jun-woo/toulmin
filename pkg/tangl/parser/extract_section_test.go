//ff:func feature=tangl type=parser control=sequence
//ff:what TestExtractSection — tests extractSection for the non-heading and heading branches
package parser

import "testing"

func TestExtractSection(t *testing.T) {
	t.Run("NotAHeading", func(t *testing.T) {
		lines := []string{"plain text", "## tangl:rules"}
		sec, next, ok := extractSection(lines, 0)
		if ok {
			t.Fatal("expected ok=false for a non-heading line")
		}
		if next != 1 {
			t.Fatalf("next = %d, want 1", next)
		}
		if sec.Name != "" {
			t.Fatalf("expected zero-value section, got %+v", sec)
		}
	})

	t.Run("Heading", func(t *testing.T) {
		lines := []string{
			"## tangl:rules",
			"line one",
			"line two",
			"## tangl:cases",
		}
		sec, next, ok := extractSection(lines, 0)
		if !ok {
			t.Fatal("expected ok=true for a tangl: heading")
		}
		if sec.Name != "rules" {
			t.Errorf("sec.Name = %q, want %q", sec.Name, "rules")
		}
		if sec.HeaderLine != 1 {
			t.Errorf("sec.HeaderLine = %d, want 1", sec.HeaderLine)
		}
		if sec.LineOffset != 2 {
			t.Errorf("sec.LineOffset = %d, want 2", sec.LineOffset)
		}
		wantLines := []string{"line one", "line two"}
		if len(sec.Lines) != len(wantLines) {
			t.Fatalf("sec.Lines = %v, want %v", sec.Lines, wantLines)
		}
		for i, l := range wantLines {
			if sec.Lines[i] != l {
				t.Errorf("sec.Lines[%d] = %q, want %q", i, sec.Lines[i], l)
			}
		}
		if next != 3 {
			t.Errorf("next = %d, want 3", next)
		}
	})
}
