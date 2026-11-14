//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what TestRenderCompareAdnKo — 관형형 leaf render: bare boolean, value/empty/negated operators, and errors
package ko

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// TestRenderCompareAdnKo checks the 관형형 surface for the bare boolean (인 / 이 아닌),
// value operators (with the equals 와/과 object particle), the value-less 비어 있는
// form, and operand sides; an unknown operator and unsupported Lhs/Rhs operands error.
func TestRenderCompareAdnKo(t *testing.T) {
	cases := []struct {
		c       ast.Compare
		negated bool
		want    string
	}{
		{ast.Compare{Field: "활성", Op: "equals", Value: "true"}, false, "`활성` 인"},
		{ast.Compare{Field: "활성", Op: "equals", Value: "true"}, true, "`활성` 이 아닌"},
		{ast.Compare{Field: "신용점수", Op: "is at least", Value: "650"}, false, "`신용점수` 가 650 이상인"},
		{ast.Compare{Field: "신용점수", Op: "is at least", Value: "650"}, true, "`신용점수` 가 650 이상이 아닌"},
		{ast.Compare{Field: "등급", Op: "equals", Value: `"VIP"`}, false, "`등급` 이 \"VIP\" 와 같은"},
		{ast.Compare{Field: "목록", Op: "is empty"}, false, "`목록` 이 비어 있는"},
		{ast.Compare{Lhs: ast.Path{Steps: []string{"c", "최소가격"}}, Op: "is at most", Rhs: ast.Path{Steps: []string{"기본가"}}}, false, "`c` 의 `최소가격` 이 `기본가` 이하인"},
	}
	for _, c := range cases {
		got, err := renderCompareAdnKo(c.c, c.negated)
		if err != nil {
			t.Errorf("renderCompareAdnKo(%+v, %v): %v", c.c, c.negated, err)
			continue
		}
		if got != c.want {
			t.Errorf("renderCompareAdnKo(%+v, %v) = %q, want %q", c.c, c.negated, got, c.want)
		}
	}
	if _, err := renderCompareAdnKo(ast.Compare{Field: "x", Op: "bogus", Value: "1"}, false); err == nil {
		t.Error("renderCompareAdnKo accepted an unknown operator")
	}
	if _, err := renderCompareAdnKo(ast.Compare{Lhs: &ast.Path{Steps: []string{"x"}}, Op: "is at least", Value: "1"}, false); err == nil {
		t.Error("renderCompareAdnKo accepted an unsupported Lhs operand")
	}
	if _, err := renderCompareAdnKo(ast.Compare{Field: "x", Op: "is at least", Rhs: &ast.Path{Steps: []string{"y"}}}, false); err == nil {
		t.Error("renderCompareAdnKo accepted an unsupported Rhs operand")
	}
}
