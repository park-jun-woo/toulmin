//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what TestRenderCompareKo — renderCompareKo covers value/empty/negated operators and the bare boolean
package ko

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// TestRenderCompareKo checks the KO predicate for value operators (with the equals
// 와/과 object particle), the value-less 비어 있음 forms, negated variants, and the bare
// boolean Compare{equals,true} → `x` 이면 / `x` 이 아니고. An unknown operator errors.
func TestRenderCompareKo(t *testing.T) {
	cases := []struct {
		c       ast.Compare
		negated bool
		want    string
	}{
		{ast.Compare{Field: "신용점수", Op: "is at least", Value: "650"}, false, "`신용점수` 가 650 이상이면"},
		{ast.Compare{Field: "등급", Op: "equals", Value: `"VIP"`}, false, "`등급` 이 \"VIP\" 와 같으면"},
		{ast.Compare{Field: "연체 이력", Op: "is empty"}, false, "`연체 이력` 이 비어 있으면"},
		{ast.Compare{Field: "x", Op: "is not empty"}, false, "`x` 가 비어 있지 않으면"},
		{ast.Compare{Field: "목록", Op: "is in", Value: "[a]"}, false, "`목록` 이 [a] 중 하나면"},
		{ast.Compare{Field: "금액", Op: "is greater than", Value: "0"}, true, "`금액` 이 0 보다 크지 않으면"},
		{ast.Compare{Field: "등급", Op: "equals", Value: `"VIP"`}, true, "`등급` 이 \"VIP\" 와 같지 않으면"},
		{ast.Compare{Field: "직원", Op: "equals", Value: "true"}, false, "`직원` 이면"},
		{ast.Compare{Field: "직원", Op: "equals", Value: "true"}, true, "`직원` 이 아니고"},
	}
	for _, c := range cases {
		got, err := renderCompareKo(c.c, c.negated)
		if err != nil {
			t.Errorf("renderCompareKo(%+v, %v): %v", c.c, c.negated, err)
			continue
		}
		if got != c.want {
			t.Errorf("renderCompareKo(%+v, %v) = %q, want %q", c.c, c.negated, got, c.want)
		}
	}
	if _, err := renderCompareKo(ast.Compare{Field: "x", Op: "bogus", Value: "1"}, false); err == nil {
		t.Error("renderCompareKo accepted an unknown operator")
	}
}
