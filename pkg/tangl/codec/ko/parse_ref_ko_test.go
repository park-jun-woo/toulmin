//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseRefKo — local name, qualified alias.name, no-backtick, dangling qualifier
package ko

import "testing"

func TestParseRefKo(t *testing.T) {
	t.Run("Local", func(t *testing.T) {
		ref, rest, ok := parseRefKo("`Allow` 실행한다")
		if !ok || ref.Alias != "" || ref.Name != "Allow" || rest != " 실행한다" {
			t.Fatalf("got %+v rest=%q ok=%v", ref, rest, ok)
		}
	})

	t.Run("Qualified", func(t *testing.T) {
		ref, rest, ok := parseRefKo("`credit`.`Threshold`")
		if !ok || ref.Alias != "credit" || ref.Name != "Threshold" || rest != "" {
			t.Fatalf("got %+v rest=%q ok=%v", ref, rest, ok)
		}
	})

	t.Run("NoBacktick", func(t *testing.T) {
		if _, _, ok := parseRefKo("Currency"); ok {
			t.Fatal("expected ok=false for a bare token")
		}
	})

	t.Run("DanglingQualifier", func(t *testing.T) {
		if _, _, ok := parseRefKo("`credit`.Threshold"); ok {
			t.Fatal("expected ok=false when the qualifier's second name is not backticked")
		}
	})
}
