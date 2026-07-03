//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseNodeModifiersKo — 0~3개 조합, 밀착 쉼표, 순서 위반 시 잔여 반환
package ko

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestParseNodeModifiersKo(t *testing.T) {
	t.Run("None", func(t *testing.T) {
		n := &ast.Node{}
		rest, err := parseNodeModifiersKo(n, "`정책`.`f` 쓰는", 1, "t.md")
		if err != nil {
			t.Fatalf("err=%v", err)
		}
		if rest != "`정책`.`f` 쓰는" || n.With != nil || n.Qualified != nil || n.Level != ast.LevelUnset {
			t.Fatalf("rest=%q n=%+v", rest, n)
		}
	})

	t.Run("WithOnly", func(t *testing.T) {
		n := &ast.Node{}
		rest, err := parseNodeModifiersKo(n, "`승인 기준` 을 기준으로 하는", 1, "t.md")
		if err != nil {
			t.Fatalf("err=%v", err)
		}
		if rest != "" || len(n.With) != 1 || n.With[0] != "승인 기준" {
			t.Fatalf("rest=%q with=%v", rest, n.With)
		}
	})

	t.Run("QualifiedOnly", func(t *testing.T) {
		n := &ast.Node{}
		rest, err := parseNodeModifiersKo(n, "가중치 0.7 의 `가격`.`Cap` 쓰는", 1, "t.md")
		if err != nil {
			t.Fatalf("err=%v", err)
		}
		if rest != "`가격`.`Cap` 쓰는" || n.Qualified == nil || *n.Qualified != 0.7 {
			t.Fatalf("rest=%q qualified=%v", rest, n.Qualified)
		}
	})

	t.Run("ReportOnly", func(t *testing.T) {
		n := &ast.Node{}
		rest, err := parseNodeModifiersKo(n, "검토를 보고하는 `심사`.`NeedsReview` 쓰는", 1, "t.md")
		if err != nil {
			t.Fatalf("err=%v", err)
		}
		if rest != "`심사`.`NeedsReview` 쓰는" || n.Level != ast.LevelReview {
			t.Fatalf("rest=%q level=%v", rest, n.Level)
		}
	})

	// All three, comma-separated (밀착 쉼표: '하는,' 과 '의,' 를 넘어 파싱).
	t.Run("AllThreeWithCommas", func(t *testing.T) {
		n := &ast.Node{}
		rest, err := parseNodeModifiersKo(n, "`차단 목록` 을 기준으로 하는, 가중치 0.8 의, 실패를 보고하는 `정책`.`IsIPBlocked` 쓰는", 1, "t.md")
		if err != nil {
			t.Fatalf("err=%v", err)
		}
		if len(n.With) != 1 || n.With[0] != "차단 목록" {
			t.Fatalf("with=%v", n.With)
		}
		if n.Qualified == nil || *n.Qualified != 0.8 {
			t.Fatalf("qualified=%v", n.Qualified)
		}
		if n.Level != ast.LevelFail {
			t.Fatalf("level=%v", n.Level)
		}
		if rest != "`정책`.`IsIPBlocked` 쓰는" {
			t.Fatalf("rest=%q", rest)
		}
	})

	// Out-of-order (가중치 before 기준) — the 기준 modifier is left unconsumed and
	// surfaces in the remainder for the caller's '쓰는' check to reject.
	t.Run("OutOfOrderLeavesRemainder", func(t *testing.T) {
		n := &ast.Node{}
		rest, err := parseNodeModifiersKo(n, "가중치 0.5 의 `a` 을 기준으로 하는", 1, "t.md")
		if err != nil {
			t.Fatalf("err=%v", err)
		}
		if n.Qualified == nil || *n.Qualified != 0.5 {
			t.Fatalf("qualified=%v", n.Qualified)
		}
		if n.With != nil {
			t.Fatalf("expected 기준 to stay unconsumed, with=%v", n.With)
		}
		if rest != "`a` 을 기준으로 하는" {
			t.Fatalf("rest=%q", rest)
		}
	})
}
