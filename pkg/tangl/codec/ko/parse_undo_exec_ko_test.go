//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseUndoExecKo — success path and every malformed-body error path
package ko

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestParseUndoExecKo(t *testing.T) {
	t.Run("Bare", func(t *testing.T) {
		ex, err := parseUndoExecKo("`잔액 충분` 면 `은행`.`refund` 로", 1, "t.md")
		if err != nil {
			t.Fatalf("err=%v", err)
		}
		if ex.Kind != ast.UndoExec || ex.Node != "잔액 충분" {
			t.Fatalf("got %+v", ex)
		}
		if ex.Func == nil || ex.Func.Alias != "은행" || ex.Func.Name != "refund" {
			t.Fatalf("func=%+v", ex.Func)
		}
	})

	t.Run("Errors", func(t *testing.T) {
		bad := []string{
			"면 `F` 로",           // missing node backtick
			"`X` 면",             // missing function ref
			"`X` 면 `F",          // unterminated function backtick (parseRefKo fails)
			"`X` 면 `F` 로 여분텍스트", // trailing text after function ref
		}
		for _, body := range bad {
			if _, err := parseUndoExecKo(body, 1, "t.md"); err == nil {
				t.Errorf("%q: expected error", body)
			}
		}
	})
}
