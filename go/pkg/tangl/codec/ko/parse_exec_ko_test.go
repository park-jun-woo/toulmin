//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseExecKo — do/undo/run dispatch, 한 번만 once, certainty gate, non-exec passthrough
package ko

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/codec/scan"
)

func TestParseExecKo(t *testing.T) {
	t.Run("DoOnce", func(t *testing.T) {
		ex, ok, err := parseExecKo(scan.Item{Text: "`주문 접수` 면 `팔`.`placeCup` 한 번만 실행한다", No: 8}, "t.md")
		if err != nil || !ok {
			t.Fatalf("ok=%v err=%v", ok, err)
		}
		if ex.Kind != ast.DoExec || !ex.Once || ex.Node != "주문 접수" || ex.Func == nil || ex.Func.Alias != "팔" || ex.Func.Name != "placeCup" {
			t.Fatalf("got %+v func=%+v", ex, ex.Func)
		}
		if ex.Certainty != nil {
			t.Fatalf("unexpected certainty: %+v", ex.Certainty)
		}
	})

	t.Run("DoNoOnce", func(t *testing.T) {
		ex, _, err := parseExecKo(scan.Item{Text: "`추출 완료` 면 `기록`.`CoffeeStatus` 실행한다", No: 1}, "t.md")
		if err != nil || ex.Once {
			t.Fatalf("got %+v err=%v", ex, err)
		}
	})

	t.Run("Certainty", func(t *testing.T) {
		ex, _, err := parseExecKo(scan.Item{Text: "`승인 가능` 확신 75% 이상이면 `승인`.`Approve` 실행한다", No: 1}, "t.md")
		if err != nil {
			t.Fatalf("err=%v", err)
		}
		if ex.Certainty == nil || ex.Certainty.Op != "at least" || ex.Certainty.Percent != 75 || ex.Once {
			t.Fatalf("got %+v cert=%+v", ex, ex.Certainty)
		}
	})

	t.Run("Undo", func(t *testing.T) {
		ex, _, err := parseExecKo(scan.Item{Text: "`잔액 충분` 이면 `은행`.`refund` 로 되돌린다", No: 1}, "t.md")
		if err != nil {
			t.Fatalf("err=%v", err)
		}
		if ex.Kind != ast.UndoExec || ex.Node != "잔액 충분" || ex.Func == nil || ex.Func.Name != "refund" {
			t.Fatalf("got %+v", ex)
		}
	})

	t.Run("Run", func(t *testing.T) {
		ex, _, err := parseExecKo(scan.Item{Text: "`주문 접수` 면 `물 붓기 가능` 진행한다", No: 1}, "t.md")
		if err != nil {
			t.Fatalf("err=%v", err)
		}
		if ex.Kind != ast.RunExec || ex.Case != "물 붓기 가능" || ex.Node != "주문 접수" || ex.Func != nil {
			t.Fatalf("got %+v", ex)
		}
	})

	t.Run("NotExec", func(t *testing.T) {
		_, ok, err := parseExecKo(scan.Item{Text: "`X` 무효다", No: 1}, "t.md")
		if err != nil || ok {
			t.Fatalf("expected ok=false passthrough, ok=%v err=%v", ok, err)
		}
	})

	t.Run("Errors", func(t *testing.T) {
		bad := []string{
			"면 `F` 실행한다",         // do: missing node backtick
			"`X` 면 실행한다",         // do: missing function ref
			"`X` 면 `F` `Z` 실행한다", // do: trailing text after ref
			"면 `F` 로 되돌린다",       // undo: missing node backtick
			"`X` 면 되돌린다",         // undo: missing function ref
			"면 `C` 진행한다",         // run: missing node backtick
			"`X` 면 진행한다",         // run: missing case name
		}
		for _, text := range bad {
			if _, ok, err := parseExecKo(scan.Item{Text: text, No: 1}, "t.md"); err == nil || !ok {
				t.Errorf("%q: expected error (ok=%v err=%v)", text, ok, err)
			}
		}
	})
}
