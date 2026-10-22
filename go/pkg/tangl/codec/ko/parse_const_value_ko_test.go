//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseConstValueKo — bare value, unit value, 타입 spec ref, missing 타입, trailing text
package ko

import "testing"

func TestParseConstValueKo(t *testing.T) {
	t.Run("BareValue", func(t *testing.T) {
		v, ref, err := parseConstValueKo("650 을", "", 1, "t.md")
		if err != nil || v != "650" || ref != nil {
			t.Fatalf("got v=%q ref=%v err=%v", v, ref, err)
		}
	})

	t.Run("QuotedString", func(t *testing.T) {
		v, ref, err := parseConstValueKo(`"go build ./..." 를`, "", 1, "t.md")
		if err != nil || v != "go build ./..." || ref != nil {
			t.Fatalf("got v=%q ref=%v err=%v", v, ref, err)
		}
	})

	t.Run("QuotedStringWithSpecRef", func(t *testing.T) {
		v, ref, err := parseConstValueKo(`"./work" 를`, ", 타입 `설정`.`Path`", 1, "t.md")
		if err != nil || v != "./work" || ref == nil {
			t.Fatalf("got v=%q ref=%v err=%v", v, ref, err)
		}
	})

	t.Run("UnmatchedQuoteIsLiteral", func(t *testing.T) {
		v, _, err := parseConstValueKo(`"미완성 을`, "", 1, "t.md")
		if err != nil || v != `"미완성` {
			t.Fatalf("got v=%q err=%v", v, err)
		}
	})

	t.Run("UnitValueWithSpecRef", func(t *testing.T) {
		v, ref, err := parseConstValueKo("65°C 를", ", 타입 `센서`.`Temperature`", 1, "t.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if v != "65°C" {
			t.Fatalf("value = %q, want 65°C", v)
		}
		if ref == nil || ref.Alias != "센서" || ref.Name != "Temperature" {
			t.Fatalf("ref = %+v", ref)
		}
	})

	t.Run("MissingTaip", func(t *testing.T) {
		if _, _, err := parseConstValueKo("650 을", ", `credit`.`Threshold`", 1, "t.md"); err == nil {
			t.Fatal("expected error for missing '타입'")
		}
	})

	t.Run("TrailingText", func(t *testing.T) {
		if _, _, err := parseConstValueKo("650 을", ", 타입 `a`.`b` `c`", 1, "t.md"); err == nil {
			t.Fatal("expected trailing-text error")
		}
	})

	t.Run("MissingRefAfterTaip", func(t *testing.T) {
		if _, _, err := parseConstValueKo("650 을", ", 타입 credit", 1, "t.md"); err == nil {
			t.Fatal("expected error for a non-backtick reference after '타입'")
		}
	})
}
