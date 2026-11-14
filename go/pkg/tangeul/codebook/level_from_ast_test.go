//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestLevelFromAST — every surface value maps to a code that maps back; an unknown value fails
package codebook

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestLevelFromAST(t *testing.T) {
	for _, v := range []ast.Level{ast.LevelUnset, ast.LevelFail, ast.LevelReview} {
		c, err := LevelFromAST(v)
		if err != nil {
			t.Fatalf("%v: %v", v, err)
		}
		if back, err := LevelToAST(c); err != nil || back != v {
			t.Fatalf("%v: round trip got %v %v", v, back, err)
		}
	}
	if _, err := LevelFromAST(ast.Level(9)); err == nil {
		t.Fatal("unknown value: want error")
	}
}
