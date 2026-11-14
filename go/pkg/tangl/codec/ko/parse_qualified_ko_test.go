//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseQualifiedKo — leading-space skip, a valid weight, the missing-digit error, invalid-value, range and precision errors
package ko

import "testing"

func TestParseQualifiedKo(t *testing.T) {
	t.Run("Ok", func(t *testing.T) {
		q, rest, err := parseQualifiedKo(" 0.5", 1, "t.md")
		if err != nil {
			t.Fatalf("err=%v", err)
		}
		if q != 0.5 || rest != "" {
			t.Fatalf("q=%v rest=%q", q, rest)
		}
	})

	t.Run("OkWithRemainder", func(t *testing.T) {
		q, rest, err := parseQualifiedKo("0.25abc", 1, "t.md")
		if err != nil {
			t.Fatalf("err=%v", err)
		}
		if q != 0.25 || rest != "abc" {
			t.Fatalf("q=%v rest=%q", q, rest)
		}
	})

	t.Run("NoDigits", func(t *testing.T) {
		if _, _, err := parseQualifiedKo("abc", 1, "t.md"); err == nil {
			t.Fatal("expected error for a missing float")
		}
	})

	t.Run("InvalidValue", func(t *testing.T) {
		if _, _, err := parseQualifiedKo("-", 1, "t.md"); err == nil {
			t.Fatal("expected error for a lone '-' that ParseFloat rejects")
		}
	})

	t.Run("OutOfRange", func(t *testing.T) {
		if _, _, err := parseQualifiedKo("3.5", 1, "t.md"); err == nil {
			t.Fatal("expected error for a weight above 1")
		}
	})

	t.Run("TooPrecise", func(t *testing.T) {
		if _, _, err := parseQualifiedKo("0.12345", 1, "t.md"); err == nil {
			t.Fatal("expected error for more than 4 decimal places")
		}
	})
}
