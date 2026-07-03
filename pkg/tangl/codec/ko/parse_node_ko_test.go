//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseNodeKo — roles×using/checking/bare, 기준 with-terms, 가중 qualified, error paths
package ko

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/pkg/tangl/codec/scan"
)

func TestParseNodeKo(t *testing.T) {
	t.Run("BareGeneral", func(t *testing.T) {
		n, ok, err := parseNodeKo(scan.Item{Text: "`인증` 은 일반 규칙이다", No: 2}, "t.md")
		if err != nil || !ok {
			t.Fatalf("ok=%v err=%v", ok, err)
		}
		if n.Name != "인증" || n.Role != ast.GeneralRule || n.Using != nil || n.Checking != "" {
			t.Fatalf("got %+v", n)
		}
	})

	t.Run("UsingCounter", func(t *testing.T) {
		n, _, err := parseNodeKo(scan.Item{Text: "`컵 없음` 은 `센서`.`noCupAvailable` 쓰는 반박 규칙이다", No: 1}, "t.md")
		if err != nil {
			t.Fatalf("err=%v", err)
		}
		if n.Role != ast.CounterRule || n.Using == nil || n.Using.Alias != "센서" || n.Using.Name != "noCupAvailable" {
			t.Fatalf("got %+v using=%+v", n, n.Using)
		}
	})

	t.Run("ExceptWithSpec", func(t *testing.T) {
		n, _, err := parseNodeKo(scan.Item{Text: "`저온` 은 `센서`.`temperatureLow` 쓰는 예외 규칙이다, 기준 `서빙 온도`", No: 1}, "t.md")
		if err != nil {
			t.Fatalf("err=%v", err)
		}
		if n.Role != ast.ExceptRule || len(n.With) != 1 || n.With[0] != "서빙 온도" {
			t.Fatalf("got %+v with=%v", n, n.With)
		}
	})

	t.Run("Checking", func(t *testing.T) {
		n, _, err := parseNodeKo(scan.Item{Text: "`합산` 은 `하위 판정` 판정 쓰는 일반 규칙이다", No: 1}, "t.md")
		if err != nil {
			t.Fatalf("err=%v", err)
		}
		if n.Checking != "하위 판정" || n.Using != nil {
			t.Fatalf("got %+v", n)
		}
	})

	t.Run("WithTermsAndQualified", func(t *testing.T) {
		n, _, err := parseNodeKo(scan.Item{Text: "`위험` 은 `정책`.`Risk` 쓰는 반박 규칙이다, 기준 `a` 그리고 `b`, 가중 0.8", No: 1}, "t.md")
		if err != nil {
			t.Fatalf("err=%v", err)
		}
		if len(n.With) != 2 || n.With[0] != "a" || n.With[1] != "b" {
			t.Fatalf("with=%v", n.With)
		}
		if n.Qualified == nil || *n.Qualified != 0.8 {
			t.Fatalf("qualified=%v", n.Qualified)
		}
	})

	t.Run("NotNode", func(t *testing.T) {
		_, ok, err := parseNodeKo(scan.Item{Text: "`사용자` 는 필수다", No: 1}, "t.md")
		if err != nil || ok {
			t.Fatalf("expected ok=false passthrough, ok=%v err=%v", ok, err)
		}
	})

	t.Run("NoBacktick", func(t *testing.T) {
		_, ok, err := parseNodeKo(scan.Item{Text: "규칙이다 없음", No: 1}, "t.md")
		if err != nil || ok {
			t.Fatalf("expected ok=false passthrough for a backtick-less line, ok=%v err=%v", ok, err)
		}
	})

	t.Run("BadRole", func(t *testing.T) {
		if _, ok, err := parseNodeKo(scan.Item{Text: "`x` 는 이상한 규칙이다", No: 1}, "t.md"); err == nil || !ok {
			t.Fatalf("expected role error (ok=%v err=%v)", ok, err)
		}
	})

	t.Run("Errors", func(t *testing.T) {
		bad := []string{
			"`x` 는 뭔가 일반 규칙이다",                  // clause not ending in 쓰는
			"`x` 는 판정 쓰는 일반 규칙이다",               // checking with no case backtick
			"`x` 는 없는함수 쓰는 일반 규칙이다",             // using with a non-backtick ref
			"`x` 는 일반 규칙이다, 알수없는 꼬리",            // unknown tail clause
			"`x` 는 `정책`.`R` 쓰는 반박 규칙이다, 가중 xyz", // bad qualified float
		}
		for _, text := range bad {
			if _, ok, err := parseNodeKo(scan.Item{Text: text, No: 1}, "t.md"); err == nil || !ok {
				t.Errorf("%q: expected error (ok=%v err=%v)", text, ok, err)
			}
		}
	})
}
