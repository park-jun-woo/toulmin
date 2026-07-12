//ff:func feature=tangl type=parser control=sequence
//ff:what renderExistsHeadKo — render a nested exists head (`` <src> 중 `c` 가 다음인 것이 있으면 ``)
package ko

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// renderExistsHeadKo renders the head line of a nested exists — one whose body is a
// multi-term condition list rather than a single leaf — as
// `<source> 중 `binder` 가 다음인 것이 있으면`. The child condition list is rendered
// separately at depth+1 by the caller (renderTermKo), exactly as the 다음 중 either
// group emits its head plus nested list.
func renderExistsHeadKo(x ast.Exists) (string, error) {
	src, _, err := renderOperandKo(x.Source)
	if err != nil {
		return "", err
	}
	return src + " 중 `" + x.Binder + "` " + josaFor(x.Binder, "이", "가") + " 다음인 것이 있으면", nil
}
