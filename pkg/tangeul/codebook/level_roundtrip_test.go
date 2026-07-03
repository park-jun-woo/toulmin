//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestLevelRoundTrip — Level↔LevelCode golden values, roundtrip identity, unknown errors
package codebook

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestLevelRoundTrip(t *testing.T) {
	cases := []struct {
		level ast.Level
		code  LevelCode
	}{
		{ast.LevelUnset, LevelUnset},
		{ast.LevelFail, LevelFail},
		{ast.LevelReview, LevelReview},
	}
	for _, c := range cases {
		code, err := LevelFromAST(c.level)
		if err != nil || code != c.code {
			t.Fatalf("LevelFromAST(%v): got %d err=%v, want %d", c.level, uint16(code), err, uint16(c.code))
		}
		level, err := LevelToAST(c.code)
		if err != nil || level != c.level {
			t.Fatalf("LevelToAST(%d): got %v err=%v, want %v", uint16(c.code), level, err, c.level)
		}
	}
	if _, err := LevelFromAST(ast.Level(99)); err == nil {
		t.Error("LevelFromAST(unknown): expected error")
	}
	if _, err := LevelToAST(LevelCode(99)); err == nil {
		t.Error("LevelToAST(unknown): expected error")
	}
}
