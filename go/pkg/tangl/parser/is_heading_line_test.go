//ff:func feature=tangl type=parser control=sequence
//ff:what TestIsHeadingLine — tests isHeadingLine for heading and non-heading branches
package parser

import "testing"

func TestIsHeadingLine(t *testing.T) {
	t.Run("True", func(t *testing.T) {
		if !isHeadingLine("  ## tangl:Subject") {
			t.Fatal("expected heading line to be detected")
		}
	})

	t.Run("False", func(t *testing.T) {
		if isHeadingLine("- not a heading") {
			t.Fatal("expected non-heading line to be false")
		}
	})
}
