//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestLevelToAST — every code maps to its surface value; an unknown code fails
package codebook

import (
	"testing"

	"github.com/geul-org/geul/go/pkg/toulmin"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestLevelToAST(t *testing.T) {
	for _, v := range []ast.Level{ast.LevelUnset, ast.LevelFail, ast.LevelReview} {
		c, _ := LevelFromAST(v)
		if got, err := LevelToAST(c); err != nil || got != v {
			t.Fatalf("%v: got %v %v", v, got, err)
		}
	}
	if _, err := LevelToAST(toulmin.LevelCode(9)); err == nil {
		t.Fatal("unknown code: want error")
	}
}
