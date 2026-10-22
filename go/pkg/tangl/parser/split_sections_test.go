//ff:func feature=tangl type=parser control=sequence
//ff:what TestSplitSections — tests splitSections for no-section, skip-non-heading, single-section, and multi-section branches
package parser

import "testing"

func TestSplitSections(t *testing.T) {
	t.Run("NoSections", func(t *testing.T) {
		secs := splitSections([]string{"just text", "more text"})
		if len(secs) != 0 {
			t.Fatalf("expected no sections, got %+v", secs)
		}
	})

	t.Run("SkipsLeadingNonHeadingLines", func(t *testing.T) {
		lines := []string{"preamble", "## tangl:Subject", "- this document is `t`"}
		secs := splitSections(lines)
		if len(secs) != 1 {
			t.Fatalf("expected 1 section, got %+v", secs)
		}
		if secs[0].Name != "Subject" {
			t.Errorf("expected Name=Subject, got %q", secs[0].Name)
		}
		if secs[0].HeaderLine != 2 {
			t.Errorf("expected HeaderLine=2, got %d", secs[0].HeaderLine)
		}
		if len(secs[0].Lines) != 1 || secs[0].Lines[0] != "- this document is `t`" {
			t.Errorf("expected Lines=[- this document is `t`], got %v", secs[0].Lines)
		}
		if secs[0].LineOffset != 3 {
			t.Errorf("expected LineOffset=3, got %d", secs[0].LineOffset)
		}
	})

	t.Run("MultipleSections", func(t *testing.T) {
		lines := []string{
			"## tangl:Subject",
			"- this document is `t`",
			"## tangl:Cases",
			"- in case of `x`",
		}
		secs := splitSections(lines)
		if len(secs) != 2 {
			t.Fatalf("expected 2 sections, got %+v", secs)
		}
		if secs[0].Name != "Subject" || len(secs[0].Lines) != 1 {
			t.Errorf("unexpected first section: %+v", secs[0])
		}
		if secs[1].Name != "Cases" || len(secs[1].Lines) != 1 || secs[1].Lines[0] != "- in case of `x`" {
			t.Errorf("unexpected second section: %+v", secs[1])
		}
	})
}
