//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseSubjectSection — tests parseSubjectSection for parseItems errors, empty items, keyword/backtick/trailing/success branches
package parser

import (
	"strings"
	"testing"
)

func TestParseSubjectSection(t *testing.T) {
	t.Run("ParseItemsError", func(t *testing.T) {
		sec := section{Lines: []string{"1. this document is `t`", "3. extra"}, LineOffset: 1}
		_, err := parseSubjectSection(sec, "test.md")
		if err == nil {
			t.Fatal("expected an error for misnumbered ordered list items")
		}
	})

	t.Run("NoItems", func(t *testing.T) {
		sec := section{Lines: nil, LineOffset: 1, HeaderLine: 1}
		_, err := parseSubjectSection(sec, "test.md")
		if err == nil {
			t.Fatal("expected an error for an empty section")
		}
		if !strings.Contains(err.Error(), "tangl:Subject requires 'this document is `name`'") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("NoKeyword", func(t *testing.T) {
		sec := section{Lines: []string{"- this thing is `t`"}, LineOffset: 1}
		_, err := parseSubjectSection(sec, "test.md")
		if err == nil {
			t.Fatal("expected an error for missing keyword")
		}
		if !strings.Contains(err.Error(), "expected 'this document is `name`'") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("NoBacktick", func(t *testing.T) {
		sec := section{Lines: []string{"- this document is t"}, LineOffset: 1}
		_, err := parseSubjectSection(sec, "test.md")
		if err == nil {
			t.Fatal("expected an error for missing backtick-quoted name")
		}
		if !strings.Contains(err.Error(), "expected backtick-quoted subject name") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("TrailingText", func(t *testing.T) {
		sec := section{Lines: []string{"- this document is `t` extra"}, LineOffset: 1}
		_, err := parseSubjectSection(sec, "test.md")
		if err == nil {
			t.Fatal("expected an error for trailing text")
		}
		if !strings.Contains(err.Error(), "unexpected trailing text") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("Success", func(t *testing.T) {
		sec := section{Lines: []string{"- this document is `t`"}, LineOffset: 1}
		name, err := parseSubjectSection(sec, "test.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if name != "t" {
			t.Errorf("expected name=t, got %q", name)
		}
	})
}
