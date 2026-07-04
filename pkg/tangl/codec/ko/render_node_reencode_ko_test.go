//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what TestRenderNodeReencodeKo — parseNodeKo(renderNodeKo(n)) ≡ n for the 7 확정 예시 seeds (Line 제외)
package ko

import (
	"reflect"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/pkg/tangl/codec/scan"
)

// TestRenderNodeReencodeKo seeds the 7 확정 정규형 예시 as ASTs (0~3 modifiers ×
// using/checking/bare), renders each to the adnominal surface, parses it back, and
// checks the round-trip is the identity on the AST (Line excluded — parse restamps it).
func TestRenderNodeReencodeKo(t *testing.T) {
	q07, q08, q09 := 0.7, 0.8, 0.9
	seeds := []ast.Node{
		{Name: "저온", Role: ast.ExceptRule, With: []string{"서빙 온도"}, Using: &ast.Ref{Alias: "센서", Name: "temperatureLow"}},
		{Name: "할인 상한", Role: ast.GeneralRule, Qualified: &q07, Using: &ast.Ref{Alias: "가격", Name: "Cap"}},
		{Name: "수동 심사", Role: ast.CounterRule, Level: ast.LevelReview, Using: &ast.Ref{Alias: "심사", Name: "NeedsReview"}},
		{Name: "IP 차단", Role: ast.CounterRule, With: []string{"차단 목록"}, Qualified: &q08, Level: ast.LevelFail, Using: &ast.Ref{Alias: "정책", Name: "IsIPBlocked"}},
		{Name: "복합 심사", Role: ast.GeneralRule, With: []string{"신용 기준", "연체 이력"}, Using: &ast.Ref{Alias: "심사", Name: "Evaluate"}},
		{Name: "중복 결제", Role: ast.CounterRule, Qualified: &q09, Checking: &ast.Ref{Name: "결제 이력"}},
		{Name: "기본 승인", Role: ast.GeneralRule, With: []string{"승인 기준"}},
	}
	for _, want := range seeds {
		surface := renderNodeKo(want)
		got, ok, err := parseNodeKo(scan.Item{Text: surface, No: 1}, "t.md")
		if err != nil || !ok {
			t.Errorf("%q: reparse ok=%v err=%v", surface, ok, err)
			continue
		}
		got.Line = 0
		if !reflect.DeepEqual(got, want) {
			t.Errorf("re-encode %q:\n got %+v\nwant %+v", surface, got, want)
		}
	}
}
