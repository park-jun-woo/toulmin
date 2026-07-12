//ff:func feature=tangl type=parser control=sequence
//ff:what TestRenderExistsKo — 관형형/참이면/exists renderers emit exactly what their parsers re-absorb
package ko

import (
	"reflect"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/pkg/tangl/codec/scan"
)

func TestRenderExistsKo(t *testing.T) {
	t.Run("CompareAdnRenderReparse", func(t *testing.T) {
		want := ast.Compare{Lhs: ast.Path{Steps: []string{"c", "최소가격"}}, Op: "is at most", Rhs: ast.Path{Steps: []string{"기본가"}}}
		text, err := renderCompareAdnKo(want, false)
		if err != nil {
			t.Fatalf("render: %v", err)
		}
		got, err := parseCompareAdnKo(text, 1, "t.md")
		if err != nil {
			t.Fatalf("reparse %q: %v", text, err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("render %q reparsed to %#v, want %#v", text, got, want)
		}
	})

	t.Run("CallPredRenderReparse", func(t *testing.T) {
		want := ast.Call{Fn: ast.Ref{Alias: "검사", Name: "isSpam"}}
		text := renderCallPredKo(want)
		got, ok, err := parseCallPredKo(text, 1, "t.md")
		if !ok || err != nil {
			t.Fatalf("reparse %q: ok=%v err=%v", text, ok, err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("render %q reparsed to %#v, want %#v", text, got, want)
		}
	})

	t.Run("ExistsOneLineRenderReparse", func(t *testing.T) {
		want := ast.Exists{
			Source: ast.Path{Steps: []string{"쿠폰들"}},
			Binder: "c",
			Cond:   ast.Compare{Lhs: ast.Path{Steps: []string{"c", "최소가격"}}, Op: "is at most", Rhs: ast.Path{Steps: []string{"기본가"}}},
		}
		text, ok, err := renderExistsPredKo(want)
		if !ok || err != nil {
			t.Fatalf("render: ok=%v err=%v", ok, err)
		}
		got, ok, err := parseExistsKo(text, nil, 1, "t.md")
		if !ok || err != nil {
			t.Fatalf("reparse %q: ok=%v err=%v", text, ok, err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("render %q reparsed to %#v, want %#v", text, got, want)
		}
	})

	t.Run("ExistsHeadRendersDaeumin", func(t *testing.T) {
		x := ast.Exists{Source: ast.Path{Steps: []string{"쿠폰들"}}, Binder: "항목"}
		head, err := renderExistsHeadKo(x)
		if err != nil {
			t.Fatalf("render: %v", err)
		}
		// 받침 binder (항목) takes 이, and the head reparses (with children) to the same source/binder.
		children := []scan.Item{{Text: "`항목` 의 `상태` 가 \"활성\" 과 같으면", No: 2}}
		got, ok, err := parseExistsKo(head, children, 1, "t.md")
		if !ok || err != nil {
			t.Fatalf("reparse %q: ok=%v err=%v", head, ok, err)
		}
		ex := got.(ast.Exists)
		if ex.Binder != "항목" || ex.Source.(ast.Path).Steps[0] != "쿠폰들" {
			t.Fatalf("got %#v", ex)
		}
	})
}
