//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseSubjectKo — subject success, missing 이다, missing prefix, trailing text
package ko

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/codec/scan"
)

func TestParseSubjectKo(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		name, err := parseSubjectKo(scan.Item{Text: "이 문서는 `아메리카노` 이다", No: 1}, "t.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if name != "아메리카노" {
			t.Fatalf("name = %q, want %q", name, "아메리카노")
		}
	})

	t.Run("MissingIda", func(t *testing.T) {
		if _, err := parseSubjectKo(scan.Item{Text: "이 문서는 `x`", No: 1}, "t.md"); err == nil {
			t.Fatal("expected error for missing 이다")
		}
	})

	t.Run("MissingPrefix", func(t *testing.T) {
		if _, err := parseSubjectKo(scan.Item{Text: "저 문서는 `x` 이다", No: 1}, "t.md"); err == nil {
			t.Fatal("expected error for missing '이 문서는' prefix")
		}
	})

	t.Run("TrailingText", func(t *testing.T) {
		_, err := parseSubjectKo(scan.Item{Text: "이 문서는 `x` `y` 이다", No: 1}, "t.md")
		if err == nil || !strings.Contains(err.Error(), "trailing") {
			t.Fatalf("expected trailing-text error, got %v", err)
		}
	})

	t.Run("MissingName", func(t *testing.T) {
		if _, err := parseSubjectKo(scan.Item{Text: "이 문서는 x 이다", No: 1}, "t.md"); err == nil {
			t.Fatal("expected error for missing backtick-quoted subject name")
		}
	})
}
