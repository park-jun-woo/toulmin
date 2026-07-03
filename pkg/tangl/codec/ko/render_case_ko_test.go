//ff:func feature=tangl type=parser control=sequence
//ff:what TestRenderCaseKo — renderCaseKo emits `` `X` 의 경우에 ``
package ko

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// TestRenderCaseKo checks the case declaration head.
func TestRenderCaseKo(t *testing.T) {
	if got := renderCaseKo(ast.Case{Name: "컵 놓기 가능"}); got != "`컵 놓기 가능` 의 경우에" {
		t.Errorf("renderCaseKo = %q", got)
	}
}
