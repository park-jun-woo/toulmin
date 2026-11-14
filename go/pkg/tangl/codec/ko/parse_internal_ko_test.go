//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseInternalHeadKo — on-event, every-tick with 까지, and unrecognized head
package ko

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/codec/scan"
)

func TestParseInternalHeadKo(t *testing.T) {
	t.Run("OnEvent", func(t *testing.T) {
		in, err := parseInternalHeadKo(scan.Item{Text: "`버튼 눌림` 발생 시", No: 1}, "t.md")
		if err != nil {
			t.Fatalf("err=%v", err)
		}
		if in.Kind != ast.OnEvent || in.Event != "버튼 눌림" {
			t.Fatalf("got %+v", in)
		}
	})

	t.Run("EveryUntil", func(t *testing.T) {
		in, err := parseInternalHeadKo(scan.Item{Text: "매 1초마다, `서빙 가능` 까지", No: 40}, "t.md")
		if err != nil {
			t.Fatalf("err=%v", err)
		}
		if in.Kind != ast.EveryTick || in.Interval != "1s" || in.Until != "서빙 가능" || in.Line != 40 {
			t.Fatalf("got %+v", in)
		}
	})

	t.Run("EveryNoUntil", func(t *testing.T) {
		in, err := parseInternalHeadKo(scan.Item{Text: "매 15분마다", No: 1}, "t.md")
		if err != nil || in.Interval != "15m" || in.Until != "" {
			t.Fatalf("got %+v err=%v", in, err)
		}
	})

	t.Run("Unrecognized", func(t *testing.T) {
		if _, err := parseInternalHeadKo(scan.Item{Text: "`X` 진행한다", No: 1}, "t.md"); err == nil {
			t.Fatal("expected error for a head that is neither 발생 시 nor 매")
		}
	})

	t.Run("OnEventNoBacktick", func(t *testing.T) {
		if _, err := parseInternalHeadKo(scan.Item{Text: "이벤트 발생 시", No: 1}, "t.md"); err == nil {
			t.Fatal("expected error when the event name before 발생 시 is not backtick-quoted")
		}
	})

	t.Run("OnEventTrailingText", func(t *testing.T) {
		if _, err := parseInternalHeadKo(scan.Item{Text: "`이벤트` 여분텍스트 발생 시", No: 1}, "t.md"); err == nil {
			t.Fatal("expected error for trailing text after the event name")
		}
	})

	t.Run("EveryBadInterval", func(t *testing.T) {
		if _, err := parseInternalHeadKo(scan.Item{Text: "매 마다", No: 1}, "t.md"); err == nil {
			t.Fatal("expected the every-clause interval error to propagate")
		}
	})
}
