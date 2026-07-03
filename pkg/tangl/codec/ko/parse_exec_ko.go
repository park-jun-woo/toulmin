//ff:func feature=tangl type=parser control=sequence
//ff:what parseExecKo — dispatch a case-body exec item by its sentence-final keyword
package ko

import (
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/pkg/tangl/codec/scan"
)

// parseExecKo dispatches an execution edge by its SOV sentence-final keyword:
// 실행한다 (do — parseDoExecKo, which also handles 한 번만 and 확신 gates), 되돌린다
// (undo — parseUndoExecKo), or 진행한다 (run — parseRunExecKo). It returns ok=false
// (no error) when the item ends in none of the three, so the caller can report an
// unrecognized case statement.
func parseExecKo(it scan.Item, path string) (ast.Exec, bool, error) {
	if body, ok := takeSuffixKo(it.Text, "실행한다"); ok {
		ex, err := parseDoExecKo(body, it.No, path)
		return ex, true, err
	}
	if body, ok := takeSuffixKo(it.Text, "되돌린다"); ok {
		ex, err := parseUndoExecKo(body, it.No, path)
		return ex, true, err
	}
	if body, ok := takeSuffixKo(it.Text, "진행한다"); ok {
		ex, err := parseRunExecKo(body, it.No, path)
		return ex, true, err
	}
	return ast.Exec{}, false, nil
}
