//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseTypeRefKo — bare built-in type, qualified ref stringified, empty input
package ko

import "testing"

func TestParseTypeRefKo(t *testing.T) {
	t.Run("BareBuiltin", func(t *testing.T) {
		typ, rest, ok := parseTypeRefKo("Currency")
		if !ok || typ != "Currency" || rest != "" {
			t.Fatalf("got %q rest=%q ok=%v", typ, rest, ok)
		}
	})

	t.Run("QualifiedRef", func(t *testing.T) {
		typ, rest, ok := parseTypeRefKo("`센서`.`Temperature`")
		if !ok || typ != "센서.Temperature" || rest != "" {
			t.Fatalf("got %q rest=%q ok=%v", typ, rest, ok)
		}
	})

	t.Run("LocalRef", func(t *testing.T) {
		typ, _, ok := parseTypeRefKo("`Threshold`")
		if !ok || typ != "Threshold" {
			t.Fatalf("got %q ok=%v", typ, ok)
		}
	})

	t.Run("Empty", func(t *testing.T) {
		if _, _, ok := parseTypeRefKo("   "); ok {
			t.Fatal("expected ok=false for empty input")
		}
	})

	t.Run("UnterminatedBacktick", func(t *testing.T) {
		if _, _, ok := parseTypeRefKo("`unterminated"); ok {
			t.Fatal("expected ok=false for an unterminated backtick-quoted ref")
		}
	})
}
