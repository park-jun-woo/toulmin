//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseFieldKo — bare built-in type, qualified ref type, missing name, missing 가진다, missing 타입
package ko

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/codec/scan"
)

func TestParseFieldKo(t *testing.T) {
	t.Run("BuiltinType", func(t *testing.T) {
		f, err := parseFieldKo(scan.Item{Text: "`연소득` 을 가진다, 타입 Currency", No: 2}, "t.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if f.Name != "연소득" || f.Type != "Currency" || f.Line != 2 {
			t.Fatalf("got %+v", f)
		}
	})

	t.Run("QualifiedRefType", func(t *testing.T) {
		f, err := parseFieldKo(scan.Item{Text: "`온도` 를 가진다, 타입 `센서`.`Temperature`", No: 3}, "t.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if f.Type != "센서.Temperature" {
			t.Fatalf("type = %q", f.Type)
		}
	})

	t.Run("MissingName", func(t *testing.T) {
		if _, err := parseFieldKo(scan.Item{Text: "연소득 을 가진다, 타입 Currency", No: 2}, "t.md"); err == nil {
			t.Fatal("expected error for missing backtick name")
		}
	})

	t.Run("MissingHas", func(t *testing.T) {
		if _, err := parseFieldKo(scan.Item{Text: "`연소득` 은 필수다", No: 2}, "t.md"); err == nil {
			t.Fatal("expected error for missing 가진다")
		}
	})

	t.Run("MissingType", func(t *testing.T) {
		if _, err := parseFieldKo(scan.Item{Text: "`연소득` 을 가진다, Currency", No: 2}, "t.md"); err == nil {
			t.Fatal("expected error for missing 타입")
		}
	})

	t.Run("TextBeforeHas", func(t *testing.T) {
		if _, err := parseFieldKo(scan.Item{Text: "`연소득` 뭔가 가진다, 타입 Currency", No: 2}, "t.md"); err == nil {
			t.Fatal("expected error for unexpected text before 가진다")
		}
	})

	t.Run("MissingTypeValue", func(t *testing.T) {
		if _, err := parseFieldKo(scan.Item{Text: "`연소득` 을 가진다, 타입", No: 2}, "t.md"); err == nil {
			t.Fatal("expected error for a missing type value after 타입")
		}
	})

	t.Run("TrailingTextAfterType", func(t *testing.T) {
		if _, err := parseFieldKo(scan.Item{Text: "`연소득` 을 가진다, 타입 Currency 그리고", No: 2}, "t.md"); err == nil {
			t.Fatal("expected trailing-text error after the type")
		}
	})
}
