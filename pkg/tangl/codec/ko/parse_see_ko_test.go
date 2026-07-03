//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseSeeKo — see success (spaced Hangul alias), missing 참조한다, missing package backtick, trailing text
package ko

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/codec/scan"
)

func TestParseSeeKo(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		see, err := parseSeeKo(scan.Item{Text: "`팔` 은 `robot/arm` 에서 참조한다", No: 2}, "t.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if see.Alias != "팔" || see.Package != "robot/arm" || see.Line != 2 {
			t.Fatalf("got %+v", see)
		}
	})

	t.Run("MissingKeyword", func(t *testing.T) {
		if _, err := parseSeeKo(scan.Item{Text: "`팔` 은 `robot/arm` 에서", No: 2}, "t.md"); err == nil {
			t.Fatal("expected error for missing 참조한다")
		}
	})

	t.Run("MissingPackage", func(t *testing.T) {
		if _, err := parseSeeKo(scan.Item{Text: "`팔` 은 참조한다", No: 2}, "t.md"); err == nil {
			t.Fatal("expected error for missing package backtick")
		}
	})

	t.Run("TrailingText", func(t *testing.T) {
		if _, err := parseSeeKo(scan.Item{Text: "`팔` 은 `robot/arm` `x` 에서 참조한다", No: 2}, "t.md"); err == nil {
			t.Fatal("expected trailing-text error")
		}
	})

	t.Run("MissingAlias", func(t *testing.T) {
		if _, err := parseSeeKo(scan.Item{Text: "팔 은 `robot/arm` 에서 참조한다", No: 2}, "t.md"); err == nil {
			t.Fatal("expected error for missing backtick-quoted alias")
		}
	})
}
