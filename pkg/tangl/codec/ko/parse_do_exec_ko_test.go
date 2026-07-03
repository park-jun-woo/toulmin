//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseDoExecKo — bare/once/certainty success paths and every malformed-body error path
package ko

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestParseDoExecKo(t *testing.T) {
	t.Run("Bare", func(t *testing.T) {
		ex, err := parseDoExecKo("`주문 접수` 면 `팔`.`placeCup`", 1, "t.md")
		if err != nil {
			t.Fatalf("err=%v", err)
		}
		if ex.Kind != ast.DoExec || ex.Once || ex.Certainty != nil || ex.Node != "주문 접수" {
			t.Fatalf("got %+v", ex)
		}
		if ex.Func == nil || ex.Func.Alias != "팔" || ex.Func.Name != "placeCup" {
			t.Fatalf("func=%+v", ex.Func)
		}
	})

	t.Run("Once", func(t *testing.T) {
		ex, err := parseDoExecKo("`주문 접수` 면 `팔`.`placeCup` 한 번만", 1, "t.md")
		if err != nil || !ex.Once {
			t.Fatalf("got %+v err=%v", ex, err)
		}
	})

	t.Run("Certainty", func(t *testing.T) {
		ex, err := parseDoExecKo("`승인 가능` 확신 75% 이상이면 `승인`.`Approve`", 1, "t.md")
		if err != nil {
			t.Fatalf("err=%v", err)
		}
		if ex.Certainty == nil || ex.Certainty.Op != "at least" || ex.Certainty.Percent != 75 {
			t.Fatalf("got %+v cert=%+v", ex, ex.Certainty)
		}
	})

	t.Run("Errors", func(t *testing.T) {
		bad := []string{
			"면 `F`",               // missing node backtick
			"`X` 면",               // no function backtick at all (fi < 0)
			"`X` 면 `F",            // unterminated function backtick (parseRefKo fails)
			"`X` 면 `F` 여분텍스트",     // trailing text after function ref
			"`X` 확신 스물% 이상이면 `F`", // malformed certainty gate
			"`X` 뭔가이상한말 `F`",      // non-확신 condition text that is not blank
		}
		for _, body := range bad {
			if _, err := parseDoExecKo(body, 1, "t.md"); err == nil {
				t.Errorf("%q: expected error", body)
			}
		}
	})
}
