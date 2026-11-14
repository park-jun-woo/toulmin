//ff:func feature=tangl type=parser control=sequence
//ff:what TestAtoiRest — tests atoiRest for valid-numeric and invalid-numeric branches
package parser

import (
	"strings"
	"testing"
)

func TestAtoiRest(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		n, rest, err := atoiRest("75% certain", 2, "test.md", 1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if n != 75 {
			t.Fatalf("expected n=75, got %d", n)
		}
		if rest != "% certain" {
			t.Fatalf("expected rest=%q, got %q", "% certain", rest)
		}
	})

	t.Run("Error", func(t *testing.T) {
		_, _, err := atoiRest("ab% certain", 2, "test.md", 3)
		if err == nil {
			t.Fatal("expected an error for a non-numeric prefix")
		}
		if !strings.Contains(err.Error(), "invalid percent") {
			t.Errorf("expected invalid percent error, got %v", err)
		}
	})
}
