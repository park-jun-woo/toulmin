//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseSupersedeKo — superseder-first (no reversal), non-supersede passthrough, missing backticks
package ko

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/codec/scan"
)

func TestParseSupersedeKo(t *testing.T) {
	t.Run("SupersederFirstNoReversal", func(t *testing.T) {
		// `컵 없음` 이면 `주문 접수` 배제한다 → Supersede{Superseder: 컵 없음, Superseded: 주문 접수}.
		// Unlike Attack, the order is NOT reversed: the first backtick is the superseder.
		sup, ok, err := parseSupersedeKo(scan.Item{Text: "`컵 없음` 이면 `주문 접수` 배제한다", No: 7}, "t.md")
		if err != nil || !ok {
			t.Fatalf("ok=%v err=%v", ok, err)
		}
		if sup.Superseder != "컵 없음" || sup.Superseded != "주문 접수" || sup.Line != 7 {
			t.Fatalf("got %+v (want Superseder=컵 없음 Superseded=주문 접수)", sup)
		}
	})

	t.Run("NotSupersede", func(t *testing.T) {
		_, ok, err := parseSupersedeKo(scan.Item{Text: "`X` 면 `C` 진행한다", No: 1}, "t.md")
		if err != nil || ok {
			t.Fatalf("expected ok=false passthrough, ok=%v err=%v", ok, err)
		}
	})

	t.Run("MissingSecondBacktick", func(t *testing.T) {
		if _, ok, err := parseSupersedeKo(scan.Item{Text: "`컵 없음` 이면 배제한다", No: 1}, "t.md"); err == nil || !ok {
			t.Fatalf("expected error for a missing target (ok=%v err=%v)", ok, err)
		}
	})

	t.Run("MissingFirstBacktick", func(t *testing.T) {
		if _, ok, err := parseSupersedeKo(scan.Item{Text: "그냥 문장 배제한다", No: 1}, "t.md"); err == nil || !ok {
			t.Fatalf("expected error for a missing superseder (ok=%v err=%v)", ok, err)
		}
	})

	t.Run("TrailingText", func(t *testing.T) {
		if _, ok, err := parseSupersedeKo(scan.Item{Text: "`A` 면 `B` extra 배제한다", No: 1}, "t.md"); err == nil || !ok {
			t.Fatalf("expected error for unexpected trailing text (ok=%v err=%v)", ok, err)
		}
	})
}
