//ff:func feature=tangl type=parser control=sequence
//ff:what TestLastTokenKo — returns the trailing whitespace-delimited token, empty for blank input
package ko

import "testing"

func TestLastTokenKo(t *testing.T) {
	t.Run("CaseOfSuffix", func(t *testing.T) {
		if got := lastTokenKo("`X` 의 경우에"); got != "경우에" {
			t.Fatalf("got %q, want %q", got, "경우에")
		}
	})

	t.Run("SingleToken", func(t *testing.T) {
		if got := lastTokenKo("무효다"); got != "무효다" {
			t.Fatalf("got %q, want %q", got, "무효다")
		}
	})

	t.Run("Blank", func(t *testing.T) {
		if got := lastTokenKo("   "); got != "" {
			t.Fatalf("got %q, want empty", got)
		}
	})
}
