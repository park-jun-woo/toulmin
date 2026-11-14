//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseAttackKo — attacker-first order, non-attack passthrough, missing backticks
package ko

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/codec/scan"
)

func TestParseAttackKo(t *testing.T) {
	t.Run("AttackerFirstOrder", func(t *testing.T) {
		// `컵 없음` 이면 `주문 접수` 무효다 → Attack{Target: 주문 접수, Attacker: 컵 없음}.
		atk, ok, err := parseAttackKo(scan.Item{Text: "`컵 없음` 이면 `주문 접수` 무효다", No: 5}, "t.md")
		if err != nil || !ok {
			t.Fatalf("ok=%v err=%v", ok, err)
		}
		if atk.Attacker != "컵 없음" || atk.Target != "주문 접수" || atk.Line != 5 {
			t.Fatalf("got %+v (want Target=주문 접수 Attacker=컵 없음)", atk)
		}
	})

	t.Run("NotAttack", func(t *testing.T) {
		_, ok, err := parseAttackKo(scan.Item{Text: "`X` 면 `C` 진행한다", No: 1}, "t.md")
		if err != nil || ok {
			t.Fatalf("expected ok=false passthrough, ok=%v err=%v", ok, err)
		}
	})

	t.Run("MissingSecondBacktick", func(t *testing.T) {
		if _, ok, err := parseAttackKo(scan.Item{Text: "`컵 없음` 이면 무효다", No: 1}, "t.md"); err == nil || !ok {
			t.Fatalf("expected error for a missing target (ok=%v err=%v)", ok, err)
		}
	})

	t.Run("MissingFirstBacktick", func(t *testing.T) {
		if _, ok, err := parseAttackKo(scan.Item{Text: "그냥 문장 무효다", No: 1}, "t.md"); err == nil || !ok {
			t.Fatalf("expected error for a missing attacker (ok=%v err=%v)", ok, err)
		}
	})

	t.Run("TrailingText", func(t *testing.T) {
		if _, ok, err := parseAttackKo(scan.Item{Text: "`A` 면 `B` extra 무효다", No: 1}, "t.md"); err == nil || !ok {
			t.Fatalf("expected error for unexpected trailing text (ok=%v err=%v)", ok, err)
		}
	})
}
