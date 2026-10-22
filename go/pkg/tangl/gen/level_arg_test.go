//ff:func feature=tangl type=codegen control=iteration dimension=1
//ff:what TestLevelArg — tests levelArg mapping for fail, review, and unset
package gen

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestLevelArg(t *testing.T) {
	cases := []struct {
		in   ast.Level
		want string
	}{
		{ast.LevelFail, "toulmin.Fail"},
		{ast.LevelReview, "toulmin.Review"},
		{ast.LevelUnset, ""},
	}
	for _, c := range cases {
		if got := levelArg(c.in); got != c.want {
			t.Errorf("levelArg(%v) = %q, want %q", c.in, got, c.want)
		}
	}
}
