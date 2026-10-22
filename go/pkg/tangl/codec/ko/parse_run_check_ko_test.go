//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseRunCheckKo — 진행한다 run / 판정한다 check kinds, passthrough, malformed name
package ko

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/codec/scan"
)

func TestParseRunCheckKo(t *testing.T) {
	t.Run("Run", func(t *testing.T) {
		name, kind, ok, err := parseRunCheckKo(scan.Item{Text: "`컵 놓기 가능` 진행한다", No: 1}, "t.md")
		if err != nil || !ok || kind != 1 || name != "컵 놓기 가능" {
			t.Fatalf("name=%q kind=%d ok=%v err=%v", name, kind, ok, err)
		}
	})

	t.Run("Check", func(t *testing.T) {
		name, kind, ok, err := parseRunCheckKo(scan.Item{Text: "`접근 가능` 판정한다", No: 1}, "t.md")
		if err != nil || !ok || kind != 2 || name != "접근 가능" {
			t.Fatalf("name=%q kind=%d ok=%v err=%v", name, kind, ok, err)
		}
	})

	t.Run("Passthrough", func(t *testing.T) {
		_, _, ok, err := parseRunCheckKo(scan.Item{Text: "`사용자` 는 필수다", No: 1}, "t.md")
		if err != nil || ok {
			t.Fatalf("expected ok=false passthrough, ok=%v err=%v", ok, err)
		}
	})

	t.Run("MissingName", func(t *testing.T) {
		if _, _, ok, err := parseRunCheckKo(scan.Item{Text: "진행한다", No: 1}, "t.md"); err == nil || !ok {
			t.Fatalf("expected error for a missing case name (ok=%v err=%v)", ok, err)
		}
	})

	t.Run("TrailingText", func(t *testing.T) {
		if _, _, ok, err := parseRunCheckKo(scan.Item{Text: "`컵 놓기 가능` 여분텍스트 진행한다", No: 1}, "t.md"); err == nil || !ok {
			t.Fatalf("expected error for trailing text after the case name (ok=%v err=%v)", ok, err)
		}
	})
}
