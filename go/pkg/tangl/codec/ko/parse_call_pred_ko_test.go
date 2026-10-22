//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseCallPredKo — 참이면 predicate call for alias/local FUNC_REF, fall-through, and errors
package ko

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestParseCallPredKo(t *testing.T) {
	t.Run("AliasCall", func(t *testing.T) {
		expr, ok, err := parseCallPredKo("`검사`.`isSpam` 가 참이면", 1, "t.md")
		if !ok || err != nil {
			t.Fatalf("ok=%v err=%v", ok, err)
		}
		c := expr.(ast.Call)
		if c.Fn.Alias != "검사" || c.Fn.Name != "isSpam" {
			t.Fatalf("got %#v", c)
		}
	})

	t.Run("LocalCall", func(t *testing.T) {
		expr, ok, err := parseCallPredKo("`f` 가 참이면", 1, "t.md")
		if !ok || err != nil {
			t.Fatalf("ok=%v err=%v", ok, err)
		}
		c := expr.(ast.Call)
		if c.Fn.Alias != "" || c.Fn.Name != "f" {
			t.Fatalf("got %#v", c)
		}
	})

	t.Run("FallThrough", func(t *testing.T) {
		_, ok, err := parseCallPredKo("`신용점수` 가 650 이상이면", 1, "t.md")
		if ok || err != nil {
			t.Fatalf("want ok=false err=nil, got ok=%v err=%v", ok, err)
		}
	})

	t.Run("NoFuncRefError", func(t *testing.T) {
		if _, ok, err := parseCallPredKo("그것이 참이면", 1, "t.md"); !ok || err == nil {
			t.Fatalf("want committed error, got ok=%v err=%v", ok, err)
		}
	})

	t.Run("TrailingTextError", func(t *testing.T) {
		if _, ok, err := parseCallPredKo("`검사`.`isSpam` 뭔가 참이면", 1, "t.md"); !ok || err == nil {
			t.Fatalf("want committed error, got ok=%v err=%v", ok, err)
		}
	})
}
