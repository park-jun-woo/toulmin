//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseNodeModifiersErrorKo — 가중치 값 비실수, 개시어 없는 레벨어는 잔여 없이 에러
package ko

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestParseNodeModifiersErrorKo(t *testing.T) {
	// 가중치 keyword present but the value is not a float → parseQualifiedKo errors.
	t.Run("QualifiedError", func(t *testing.T) {
		n := &ast.Node{}
		rest, err := parseNodeModifiersKo(n, "가중치 없음 의", 1, "t.md")
		if err == nil {
			t.Fatalf("expected error, rest=%q n=%+v", rest, n)
		}
		if rest != "" {
			t.Fatalf("rest=%q", rest)
		}
	})

	// Level word ('실패') without a following '보고하는' → parseReportKo errors.
	t.Run("ReportError", func(t *testing.T) {
		n := &ast.Node{}
		rest, err := parseNodeModifiersKo(n, "실패를 알리는", 1, "t.md")
		if err == nil {
			t.Fatalf("expected error, rest=%q n=%+v", rest, n)
		}
		if rest != "" {
			t.Fatalf("rest=%q", rest)
		}
	})
}
