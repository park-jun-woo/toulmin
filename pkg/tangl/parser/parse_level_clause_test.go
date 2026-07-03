//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseLevelClause — tests parseLevelClause for fail, review, missing level, and bad token branches
package parser

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestParseLevelClause(t *testing.T) {
	t.Run("Fail", func(t *testing.T) {
		lvl, rest, err := parseLevelClause("level fail", "test.md", 1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if lvl != ast.LevelFail {
			t.Fatalf("lvl = %v, want LevelFail", lvl)
		}
		if strings.TrimSpace(rest) != "" {
			t.Fatalf("rest = %q, want empty", rest)
		}
	})

	t.Run("Review", func(t *testing.T) {
		lvl, _, err := parseLevelClause("level review", "test.md", 2)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if lvl != ast.LevelReview {
			t.Fatalf("lvl = %v, want LevelReview", lvl)
		}
	})

	t.Run("MissingLevelKeyword", func(t *testing.T) {
		_, _, err := parseLevelClause("fail", "test.md", 3)
		if err == nil || !strings.Contains(err.Error(), "expected 'level'") {
			t.Fatalf("expected level-keyword error, got %v", err)
		}
	})

	t.Run("BadToken", func(t *testing.T) {
		_, _, err := parseLevelClause("level foo", "test.md", 4)
		if err == nil || !strings.Contains(err.Error(), "expected 'fail' or 'review'") {
			t.Fatalf("expected bad-token error, got %v", err)
		}
	})
}
