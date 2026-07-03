//ff:func feature=tangl type=parser control=selection
//ff:what renderNodeKo — render a case node registration (`` `X` 는 [<fn>|`Y` 판정] 쓰는 <역할> 규칙이다[꼬리] ``)
package ko

import (
	"strconv"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// renderNodeKo renders a case node back to "`<name>` 는 [<using> 쓰는] <role> 규칙이다"
// plus its optional trailing clauses: a function ref '<ref> 쓰는' or verdict '`<case>` 판정 쓰는'
// before the role, then ', 기준 <terms>' (with), ', 가중 <float>' (qualified), and
// ', 수준 실패'/', 수준 검토' (level). The role word is 일반/반박/예외 (general/counter/except).
// It inverts parseNodeKo/applyNodeTailKo.
func renderNodeKo(n ast.Node) string {
	role := "일반"
	switch n.Role {
	case ast.CounterRule:
		role = "반박"
	case ast.ExceptRule:
		role = "예외"
	}
	s := "`" + n.Name + "` " + josaFor(n.Name, "은", "는") + " "
	if n.Using != nil {
		s += renderRefKo(*n.Using) + " 쓰는 "
	} else if n.Checking != "" {
		s += "`" + n.Checking + "` 판정 쓰는 "
	}
	s += role + " 규칙이다" + renderWithKo(n.With)
	if n.Qualified != nil {
		s += ", 가중 " + strconv.FormatFloat(*n.Qualified, 'g', -1, 64)
	}
	switch n.Level {
	case ast.LevelFail:
		s += ", 수준 실패"
	case ast.LevelReview:
		s += ", 수준 검토"
	}
	return s
}
