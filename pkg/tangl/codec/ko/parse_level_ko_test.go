//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseLevelKo — 실패→Fail, 검토→Review, and an unrecognized value error
package ko

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestParseLevelKo(t *testing.T) {
	t.Run("Fail", func(t *testing.T) {
		lvl, rest, err := parseLevelKo("실패", 3, "t.md")
		if err != nil || lvl != ast.LevelFail || rest != "" {
			t.Fatalf("lvl=%v rest=%q err=%v", lvl, rest, err)
		}
	})

	t.Run("Review", func(t *testing.T) {
		lvl, rest, err := parseLevelKo("검토", 4, "t.md")
		if err != nil || lvl != ast.LevelReview || rest != "" {
			t.Fatalf("lvl=%v rest=%q err=%v", lvl, rest, err)
		}
	})

	t.Run("Unknown", func(t *testing.T) {
		if _, _, err := parseLevelKo("보류", 5, "t.md"); err == nil {
			t.Fatal("expected error for an unrecognized 수준 value")
		}
	})
}
