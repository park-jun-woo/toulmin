//ff:func feature=tangl type=parser control=sequence
//ff:what TestRenderExistsPredKo — one-line exists for Compare/Not leaves, nested fallback, and error paths
package ko

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestRenderExistsPredKo(t *testing.T) {
	src := ast.Path{Steps: []string{"쿠폰들"}}

	t.Run("CompareLeaf", func(t *testing.T) {
		x := ast.Exists{Source: src, Binder: "c", Cond: ast.Compare{Field: "활성", Op: "equals", Value: "true"}}
		got, ok, err := renderExistsPredKo(x)
		if !ok || err != nil {
			t.Fatalf("ok=%v err=%v", ok, err)
		}
		if want := "`쿠폰들` 중 `c` 가 `활성` 인 것이 있으면"; got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
	})

	t.Run("NotCompareLeaf", func(t *testing.T) {
		x := ast.Exists{Source: src, Binder: "c", Cond: ast.Not{Term: ast.Compare{Field: "활성", Op: "equals", Value: "true"}}}
		got, ok, err := renderExistsPredKo(x)
		if !ok || err != nil {
			t.Fatalf("ok=%v err=%v", ok, err)
		}
		if want := "`쿠폰들` 중 `c` 가 `활성` 이 아닌 것이 있으면"; got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
	})

	t.Run("NotNonCompareFallsBack", func(t *testing.T) {
		x := ast.Exists{Source: src, Binder: "c", Cond: ast.Not{Term: ast.Logic{Op: "and"}}}
		if _, ok, err := renderExistsPredKo(x); ok || err != nil {
			t.Fatalf("want ok=false err=nil, got ok=%v err=%v", ok, err)
		}
	})

	t.Run("LogicBodyFallsBack", func(t *testing.T) {
		x := ast.Exists{Source: src, Binder: "c", Cond: ast.Logic{Op: "and"}}
		if _, ok, err := renderExistsPredKo(x); ok || err != nil {
			t.Fatalf("want ok=false err=nil, got ok=%v err=%v", ok, err)
		}
	})

	t.Run("InnerRenderError", func(t *testing.T) {
		x := ast.Exists{Source: src, Binder: "c", Cond: ast.Compare{Field: "x", Op: "bogus", Value: "1"}}
		if _, _, err := renderExistsPredKo(x); err == nil {
			t.Fatal("expected error for an unknown inner operator")
		}
	})

	t.Run("SourceError", func(t *testing.T) {
		x := ast.Exists{Source: &ast.Path{Steps: []string{"쿠폰들"}}, Binder: "c", Cond: ast.Compare{Field: "활성", Op: "equals", Value: "true"}}
		if _, _, err := renderExistsPredKo(x); err == nil {
			t.Fatal("expected error for an unsupported source operand")
		}
	})
}
