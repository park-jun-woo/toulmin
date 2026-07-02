//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseItemLine — tests parseItemLine for unordered, ordered, malformed, and overflow branches
package parser

import "testing"

func TestParseItemLine(t *testing.T) {
	t.Run("Unordered", func(t *testing.T) {
		content, ordered, number, ok := parseItemLine("- hello world")
		if !ok {
			t.Fatal("expected ok=true")
		}
		if ordered {
			t.Error("expected ordered=false")
		}
		if number != 0 {
			t.Errorf("expected number=0, got %d", number)
		}
		if content != "hello world" {
			t.Errorf("expected content='hello world', got %q", content)
		}
	})

	t.Run("NoDigitsNoPrefix", func(t *testing.T) {
		_, _, _, ok := parseItemLine("hello world")
		if ok {
			t.Fatal("expected ok=false for text with no list marker")
		}
	})

	t.Run("DigitsNoTerminator", func(t *testing.T) {
		_, _, _, ok := parseItemLine("12")
		if ok {
			t.Fatal("expected ok=false for digits with no trailing '. '")
		}
	})

	t.Run("DigitsWrongPunct", func(t *testing.T) {
		_, _, _, ok := parseItemLine("12x foo")
		if ok {
			t.Fatal("expected ok=false when the char after digits is not '.'")
		}
	})

	t.Run("DigitsNoSpaceAfterDot", func(t *testing.T) {
		_, _, _, ok := parseItemLine("12.foo")
		if ok {
			t.Fatal("expected ok=false when there is no space after '.'")
		}
	})

	t.Run("AtoiOverflow", func(t *testing.T) {
		_, _, _, ok := parseItemLine("123456789012345678901234567890. foo")
		if ok {
			t.Fatal("expected ok=false for an integer overflow")
		}
	})

	t.Run("OrderedSuccess", func(t *testing.T) {
		content, ordered, number, ok := parseItemLine("12. item text")
		if !ok {
			t.Fatal("expected ok=true")
		}
		if !ordered {
			t.Error("expected ordered=true")
		}
		if number != 12 {
			t.Errorf("expected number=12, got %d", number)
		}
		if content != "item text" {
			t.Errorf("expected content='item text', got %q", content)
		}
	})
}
