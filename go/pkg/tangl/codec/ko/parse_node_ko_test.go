//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseNodeKo — 확정 예시 7종, bare, 오류 경로, 구 꼬리절 파싱 불가 회귀
package ko

import (
	"reflect"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/codec/scan"
)

func TestParseNodeKo(t *testing.T) {
	q07, q08, q09 := 0.7, 0.8, 0.9
	// 확정 정규형 예시 7종 (TANGEUL.ko.md §사안) + bare general, each → its expected Node.
	good := []struct {
		text string
		want ast.Node
	}{
		{"`저온` 은 `서빙 온도` 를 기준으로 하는 `센서`.`temperatureLow` 쓰는 예외 규칙이다", ast.Node{Name: "저온", Role: ast.ExceptRule, With: []string{"서빙 온도"}, Using: &ast.Ref{Alias: "센서", Name: "temperatureLow"}, Line: 1}},
		{"`할인 상한` 은 가중치 0.7 의 `가격`.`Cap` 쓰는 일반 규칙이다", ast.Node{Name: "할인 상한", Role: ast.GeneralRule, Qualified: &q07, Using: &ast.Ref{Alias: "가격", Name: "Cap"}, Line: 1}},
		{"`수동 심사` 는 검토를 보고하는 `심사`.`NeedsReview` 쓰는 반박 규칙이다", ast.Node{Name: "수동 심사", Role: ast.CounterRule, Level: ast.LevelReview, Using: &ast.Ref{Alias: "심사", Name: "NeedsReview"}, Line: 1}},
		{"`IP 차단` 은 `차단 목록` 을 기준으로 하는, 가중치 0.8 의, 실패를 보고하는 `정책`.`IsIPBlocked` 쓰는 반박 규칙이다", ast.Node{Name: "IP 차단", Role: ast.CounterRule, With: []string{"차단 목록"}, Qualified: &q08, Level: ast.LevelFail, Using: &ast.Ref{Alias: "정책", Name: "IsIPBlocked"}, Line: 1}},
		{"`복합 심사` 는 `신용 기준` 그리고 `연체 이력` 을 기준으로 하는 `심사`.`Evaluate` 쓰는 일반 규칙이다", ast.Node{Name: "복합 심사", Role: ast.GeneralRule, With: []string{"신용 기준", "연체 이력"}, Using: &ast.Ref{Alias: "심사", Name: "Evaluate"}, Line: 1}},
		{"`중복 결제` 는 가중치 0.9 의 `결제 이력` 판정 쓰는 반박 규칙이다", ast.Node{Name: "중복 결제", Role: ast.CounterRule, Qualified: &q09, Checking: &ast.Ref{Name: "결제 이력"}, Line: 1}},
		{"`기본 승인` 은 `승인 기준` 을 기준으로 하는 일반 규칙이다", ast.Node{Name: "기본 승인", Role: ast.GeneralRule, With: []string{"승인 기준"}, Line: 1}},
		{"`인증` 은 일반 규칙이다", ast.Node{Name: "인증", Role: ast.GeneralRule, Line: 1}},
	}
	t.Run("Good", func(t *testing.T) {
		for _, c := range good {
			n, ok, err := parseNodeKo(scan.Item{Text: c.text, No: 1}, "t.md")
			if err != nil || !ok || !reflect.DeepEqual(n, c.want) {
				t.Errorf("%q: ok=%v err=%v\n got %+v\nwant %+v", c.text, ok, err, n, c.want)
			}
		}
	})

	// Passthrough (ok=false, no error): non-node lines, and any line where 규칙이다 is
	// not sentence-final (the retired tail forms are one such case — trailing text after
	// 규칙이다 makes the suffix dispatch decline). In a real Cases body these surface as
	// "unrecognized case statement" at applyCaseChildKo.
	t.Run("Passthrough", func(t *testing.T) {
		pass := []string{
			"`사용자` 는 필수다",
			"규칙이다 없음",
			"`저온` 은 `센서`.`temperatureLow` 쓰는 예외 규칙이다 그리고 추가",
			"`수동 심사` 는 `심사`.`NeedsReview` 쓰는 반박 규칙이다 잔여",
		}
		for _, text := range pass {
			if _, ok, err := parseNodeKo(scan.Item{Text: text, No: 1}, "t.md"); err != nil || ok {
				t.Errorf("%q: expected ok=false passthrough, got ok=%v err=%v", text, ok, err)
			}
		}
	})

	// Malformed node declarations (ok=true, error).
	t.Run("Errors", func(t *testing.T) {
		bad := []string{
			"`x` 는 이상한 규칙이다",                      // bad role word
			"`x` 는 뭔가 일반 규칙이다",                    // modifier junk not ending in 쓰는
			"`x` 는 판정 쓰는 일반 규칙이다",                 // checking with no case backtick
			"`x` 는 없는함수 쓰는 일반 규칙이다",               // using with a non-backtick ref
			"`x` 는 가중치 xyz 의 `정책`.`R` 쓰는 반박 규칙이다", // bad qualified float
			"`x` 는 실패를 쓰는 일반 규칙이다",                // level word without 보고하는
		}
		for _, text := range bad {
			if _, ok, err := parseNodeKo(scan.Item{Text: text, No: 1}, "t.md"); err == nil || !ok {
				t.Errorf("%q: expected error (ok=%v err=%v)", text, ok, err)
			}
		}
	})
}
