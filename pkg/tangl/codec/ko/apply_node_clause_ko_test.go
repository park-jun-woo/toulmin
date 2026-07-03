//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what TestApplyNodeClauseKo — 기준/가중 clauses, their error paths, and the unrecognized default
package ko

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestApplyNodeClauseKo(t *testing.T) {
	cases := []struct {
		name    string
		text    string
		wantErr bool
	}{
		{"WithTermsOk", "기준 `a`", false},
		{"WithTermsErr", "기준 노백틱", true},
		{"QualifiedOk", "가중 0.7", false},
		{"QualifiedErr", "가중 abc", true},
		{"Unrecognized", "알수없음", true},
	}
	for _, c := range cases {
		n := &ast.Node{}
		rest, err := applyNodeClauseKo(n, c.text, 1, "t.md")
		gotErr := err != nil
		if gotErr != c.wantErr {
			t.Errorf("%s: err=%v (rest=%q), wantErr=%v", c.name, err, rest, c.wantErr)
		}
	}

	t.Run("WithFilled", func(t *testing.T) {
		n := &ast.Node{}
		rest, err := applyNodeClauseKo(n, "기준 `a`", 1, "t.md")
		if err != nil {
			t.Fatalf("err=%v", err)
		}
		if rest != "" || len(n.With) != 1 || n.With[0] != "a" {
			t.Fatalf("rest=%q with=%v", rest, n.With)
		}
	})

	t.Run("QualifiedFilled", func(t *testing.T) {
		n := &ast.Node{}
		rest, err := applyNodeClauseKo(n, "가중 0.7", 1, "t.md")
		if err != nil {
			t.Fatalf("err=%v", err)
		}
		if rest != "" || n.Qualified == nil || *n.Qualified != 0.7 {
			t.Fatalf("rest=%q qualified=%v", rest, n.Qualified)
		}
	})
}
