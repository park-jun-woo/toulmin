//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what TestApplyNodeTailKo — empty tail, single/multiple 기준·가중 clauses, and a clause error
package ko

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestApplyNodeTailKo(t *testing.T) {
	cases := []struct {
		name    string
		tail    string
		wantErr bool
	}{
		{"Empty", "", false},
		{"SingleClause", "기준 `a`", false},
		{"MultipleClauses", "기준 `a`, 가중 0.5", false},
		{"ClauseError", "알수없는꼬리", true},
	}
	for _, c := range cases {
		n := &ast.Node{}
		err := applyNodeTailKo(n, c.tail, 1, "t.md")
		gotErr := err != nil
		if gotErr != c.wantErr {
			t.Errorf("%s: err=%v, wantErr=%v", c.name, err, c.wantErr)
		}
	}

	t.Run("MultipleClausesFilled", func(t *testing.T) {
		n := &ast.Node{}
		if err := applyNodeTailKo(n, "기준 `a`, 가중 0.5", 1, "t.md"); err != nil {
			t.Fatalf("err=%v", err)
		}
		if len(n.With) != 1 || n.With[0] != "a" || n.Qualified == nil || *n.Qualified != 0.5 {
			t.Fatalf("got with=%v qualified=%v", n.With, n.Qualified)
		}
	})
}
