//ff:func feature=tangl type=parser control=sequence
//ff:what TestRenderExprKo — the new-syntax renderers emit exactly what their parsers re-absorb (render→reparse identity)
package ko

import (
	"reflect"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestRenderExprKo(t *testing.T) {
	t.Run("Operand", func(t *testing.T) {
		cases := []struct {
			op       ast.Operand
			want     string
			wantBase string
		}{
			{ast.Path{Steps: []string{"작성자", "신뢰점수"}}, "`작성자` 의 `신뢰점수`", "신뢰점수"},
			{ast.SpecRef{Steps: []string{"할인", "최소 가격"}}, "기준 `할인` 의 `최소 가격`", "최소 가격"},
			{ast.Call{Fn: ast.Ref{Alias: "nlp", Name: "score"}}, "`nlp`.`score` 의 결과", "결과"},
		}
		for _, c := range cases {
			text, base, err := renderOperandKo(c.op)
			if err != nil || text != c.want || base != c.wantBase {
				t.Fatalf("%#v: got (%q,%q,%v)", c.op, text, base, err)
			}
		}
	})

	t.Run("CompareRenderReparse", func(t *testing.T) {
		cmps := []ast.Compare{
			{Lhs: ast.Path{Steps: []string{"작성자", "신뢰점수"}}, Op: "is at least", Value: "0.7"},
			{Field: "기본가", Op: "is at least", Rhs: ast.SpecRef{Steps: []string{"할인", "최소 가격"}}},
			{Field: "등급", Op: "equals", Rhs: ast.Path{Steps: []string{"작성자", "등급"}}},
			{Lhs: ast.Call{Fn: ast.Ref{Alias: "nlp", Name: "score"}}, Op: "is greater than", Value: "0.7"},
		}
		for _, want := range cmps {
			text, err := renderCompareKo(want, false)
			if err != nil {
				t.Fatalf("%#v: render: %v", want, err)
			}
			got, err := parseCompareKo(text, 1, "t.md")
			if err != nil {
				t.Fatalf("%q: reparse: %v", text, err)
			}
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("render %q reparsed to %#v, want %#v", text, got, want)
			}
		}
	})

	t.Run("BareBoolPathRenderReparse", func(t *testing.T) {
		// A bare-boolean path Lhs renders `x` 이면 and is re-absorbed by the cond-term dispatch.
		want := ast.Compare{Lhs: ast.Path{Steps: []string{"작성자", "활성"}}, Op: "equals", Value: "true"}
		text, err := renderCompareKo(want, false)
		if err != nil {
			t.Fatalf("render: %v", err)
		}
		got, err := parseCondTermKo(text, nil, 1, "t.md")
		if err != nil {
			t.Fatalf("reparse %q: %v", text, err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("render %q reparsed to %#v, want %#v", text, got, want)
		}
	})

}
