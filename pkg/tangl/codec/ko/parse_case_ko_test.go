//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseCaseKo — case head name, missing 경우에, missing backtick, trailing text
package ko

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/codec/scan"
)

func TestParseCaseKo(t *testing.T) {
	t.Run("Head", func(t *testing.T) {
		c, err := parseCaseKo(scan.Item{Text: "`컵 놓기 가능` 의 경우에", No: 7}, "t.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if c.Name != "컵 놓기 가능" || c.Line != 7 {
			t.Fatalf("got %+v", c)
		}
	})

	t.Run("MissingKeyword", func(t *testing.T) {
		if _, err := parseCaseKo(scan.Item{Text: "`X` 무효다", No: 1}, "t.md"); err == nil {
			t.Fatal("expected error when 경우에 is absent")
		}
	})

	t.Run("MissingBacktick", func(t *testing.T) {
		if _, err := parseCaseKo(scan.Item{Text: "컵 의 경우에", No: 1}, "t.md"); err == nil {
			t.Fatal("expected error for missing backtick name")
		}
	})

	t.Run("TrailingText", func(t *testing.T) {
		if _, err := parseCaseKo(scan.Item{Text: "`X` `Y` 의 경우에", No: 1}, "t.md"); err == nil {
			t.Fatal("expected error for unexpected trailing text")
		}
	})
}
