//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseRunExecKo — success path and every malformed-body error path
package ko

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestParseRunExecKo(t *testing.T) {
	t.Run("Bare", func(t *testing.T) {
		ex, err := parseRunExecKo("`주문 접수` 면 `물 붓기 가능`", 1, "t.md")
		if err != nil {
			t.Fatalf("err=%v", err)
		}
		if ex.Kind != ast.RunExec || ex.Node != "주문 접수" || ex.Case != "물 붓기 가능" || ex.Func != nil {
			t.Fatalf("got %+v", ex)
		}
	})

	t.Run("Errors", func(t *testing.T) {
		bad := []string{
			"면 `C`",           // missing node backtick
			"`X` 면",           // missing case backtick
			"`X` 면 `C` 여분텍스트", // trailing text after case name
		}
		for _, body := range bad {
			if _, err := parseRunExecKo(body, 1, "t.md"); err == nil {
				t.Errorf("%q: expected error", body)
			}
		}
	})
}
