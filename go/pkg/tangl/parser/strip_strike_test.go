//ff:func feature=tangl type=parser control=sequence
//ff:what TestStripStrike — tests stripStrike for no-prefix, no-suffix, too-short, and struck branches
package parser

import "testing"

func TestStripStrike(t *testing.T) {
	t.Run("NoPrefix", func(t *testing.T) {
		text, struck := stripStrike("abc~~")
		if struck {
			t.Fatal("expected struck=false without a '~~' prefix")
		}
		if text != "abc~~" {
			t.Errorf("expected text='abc~~', got %q", text)
		}
	})

	t.Run("NoSuffix", func(t *testing.T) {
		text, struck := stripStrike("~~abc")
		if struck {
			t.Fatal("expected struck=false without a '~~' suffix")
		}
		if text != "~~abc" {
			t.Errorf("expected text='~~abc', got %q", text)
		}
	})

	t.Run("TooShort", func(t *testing.T) {
		text, struck := stripStrike("~~")
		if struck {
			t.Fatal("expected struck=false when content is shorter than 4 runes")
		}
		if text != "~~" {
			t.Errorf("expected text='~~', got %q", text)
		}
	})

	t.Run("Struck", func(t *testing.T) {
		text, struck := stripStrike("~~abc~~")
		if !struck {
			t.Fatal("expected struck=true")
		}
		if text != "abc" {
			t.Errorf("expected text='abc', got %q", text)
		}
	})
}
