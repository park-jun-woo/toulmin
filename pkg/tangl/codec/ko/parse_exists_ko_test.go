//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseExistsKo — one-line and 다음인 nested exists, lookahead fall-through, and tail/child errors
package ko

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/pkg/tangl/codec/scan"
)

func TestParseExistsKo(t *testing.T) {
	t.Run("OneLine", func(t *testing.T) {
		expr, ok, err := parseExistsKo("`쿠폰들` 중 `c` 가 `c` 의 `최소가격` 이 `기본가` 이하인 것이 있으면", nil, 1, "t.md")
		if !ok || err != nil {
			t.Fatalf("ok=%v err=%v", ok, err)
		}
		ex := expr.(ast.Exists)
		if ex.Binder != "c" {
			t.Fatalf("binder = %q", ex.Binder)
		}
		if src := ex.Source.(ast.Path); src.Steps[0] != "쿠폰들" {
			t.Fatalf("source = %#v", ex.Source)
		}
		if _, ok := ex.Cond.(ast.Compare); !ok {
			t.Fatalf("cond = %#v, want Compare", ex.Cond)
		}
	})

	t.Run("Nested", func(t *testing.T) {
		children := []scan.Item{
			{Text: "`c` 의 `상태` 가 \"활성\" 과 같으면", No: 2},
			{Text: "그리고 `c` 의 `최소가격` 이 `기본가` 보다 작으면", No: 3},
		}
		expr, ok, err := parseExistsKo("`쿠폰들` 중 `c` 가 다음인 것이 있으면", children, 1, "t.md")
		if !ok || err != nil {
			t.Fatalf("ok=%v err=%v", ok, err)
		}
		ex := expr.(ast.Exists)
		if _, ok := ex.Cond.(ast.Logic); !ok {
			t.Fatalf("cond = %#v, want Logic", ex.Cond)
		}
	})

	t.Run("FallThroughNoJung", func(t *testing.T) {
		// `f` 가 목록 중 하나면 — operand `f` is followed by 조사 가, not a bare 중; not exists.
		_, ok, err := parseExistsKo("`f` 가 목록 중 하나면", nil, 1, "t.md")
		if ok || err != nil {
			t.Fatalf("want ok=false err=nil, got ok=%v err=%v", ok, err)
		}
	})

	t.Run("FallThroughNoBinder", func(t *testing.T) {
		// `다음 중` group head — 중 not followed by a backtick binder.
		_, ok, err := parseExistsKo("`목록` 중 하나면", nil, 1, "t.md")
		if ok || err != nil {
			t.Fatalf("want ok=false err=nil, got ok=%v err=%v", ok, err)
		}
	})

	t.Run("FallThroughBadSource", func(t *testing.T) {
		// A bare (non-backtick) source is not an operand; the lookahead falls through.
		_, ok, err := parseExistsKo("목록 중 `c` 가 다음인 것이 있으면", nil, 1, "t.md")
		if ok || err != nil {
			t.Fatalf("want ok=false err=nil, got ok=%v err=%v", ok, err)
		}
	})

	t.Run("AdnBodyError", func(t *testing.T) {
		_, ok, err := parseExistsKo("`목록` 중 `항목` 이 웅앵웅 것이 있으면", nil, 1, "t.md")
		if !ok || err == nil {
			t.Fatalf("want committed error from the adnominal body, got ok=%v err=%v", ok, err)
		}
	})

	t.Run("MissingTail", func(t *testing.T) {
		_, ok, err := parseExistsKo("`쿠폰들` 중 `c` 가 `c` 의 `최소가격` 이 `기본가` 이하인", nil, 1, "t.md")
		if !ok || err == nil {
			t.Fatalf("want committed error, got ok=%v err=%v", ok, err)
		}
	})

	t.Run("NestedNoChildren", func(t *testing.T) {
		_, ok, err := parseExistsKo("`쿠폰들` 중 `c` 가 다음인 것이 있으면", nil, 1, "t.md")
		if !ok || err == nil {
			t.Fatalf("want committed error from empty child list, got ok=%v err=%v", ok, err)
		}
	})
}
