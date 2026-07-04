//ff:func feature=tangl type=parser control=selection
//ff:what renderNodeKo — render a node (`` `X` 은 [관형절들] [<fn>|`Y` 판정 쓰는] <역할> 규칙이다 ``)
package ko

import (
	"strconv"
	"strings"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// renderNodeKo renders a case node in the adnominal canonical form: the subject,
// then the fixed-order modifiers — '<terms> 을 기준으로 하는' (With), '가중치 <float> 의'
// (Qualified), '실패를/검토를 보고하는' (Level) — joined by ", " when two or more,
// then the optional '<ref> 쓰는'/'`Y` 판정 쓰는' clause, the role word, and 규칙이다.
// It inverts parseNodeKo/parseNodeModifiersKo.
func renderNodeKo(n ast.Node) string {
	role := "일반"
	switch n.Role {
	case ast.CounterRule:
		role = "반박"
	case ast.ExceptRule:
		role = "예외"
	}
	var mods []string
	if len(n.With) > 0 {
		mods = append(mods, renderWithKo(n.With))
	}
	if n.Qualified != nil {
		mods = append(mods, "가중치 "+strconv.FormatFloat(*n.Qualified, 'g', -1, 64)+" 의")
	}
	switch n.Level {
	case ast.LevelFail:
		mods = append(mods, "실패를 보고하는")
	case ast.LevelReview:
		mods = append(mods, "검토를 보고하는")
	}
	s := "`" + n.Name + "` " + josaFor(n.Name, "은", "는") + " "
	if len(mods) > 0 {
		s += strings.Join(mods, ", ") + " "
	}
	if n.Using != nil {
		s += renderRefKo(*n.Using) + " 쓰는 "
	} else if n.Checking != nil {
		s += "`" + n.Checking.Name + "` 판정 쓰는 "
	}
	return s + role + " 규칙이다"
}
