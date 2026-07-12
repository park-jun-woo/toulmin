//ff:func feature=tangl type=parser control=selection
//ff:what renderExistsPredKo — render a one-line exists (`` <src> 중 `c` 가 <관형형> 것이 있으면 ``)
package ko

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// renderExistsPredKo renders an Exists whose body is a single leaf — a Compare or a
// Not{Compare} — as the one-line form `<source> 중 `binder` 가 <관형형 조건> 것이 있으면`
// (the source carries no 조사 before 중, mirroring the value-중 convention). ok is false
// for any other Cond (a Logic/Either/nested Exists body), so the caller falls back to
// the nested 다음인 head form. The inner condition uses the 관형형 renderer.
func renderExistsPredKo(x ast.Exists) (string, bool, error) {
	var inner string
	var err error
	switch cond := x.Cond.(type) {
	case ast.Compare:
		inner, err = renderCompareAdnKo(cond, false)
	case ast.Not:
		cmp, ok := cond.Term.(ast.Compare)
		if !ok {
			return "", false, nil
		}
		inner, err = renderCompareAdnKo(cmp, true)
	default:
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	src, _, err := renderOperandKo(x.Source)
	if err != nil {
		return "", false, err
	}
	return src + " 중 `" + x.Binder + "` " + josaFor(x.Binder, "이", "가") + " " + inner + " 것이 있으면", true, nil
}
