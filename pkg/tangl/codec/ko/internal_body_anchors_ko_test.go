//ff:func feature=tangl type=parser control=sequence
//ff:what TestInternalBodyAnchorsKo — struck skip, run/check dispatch, malformed and unrecognized errors
package ko

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/pkg/tangl/codec/scan"
)

func TestInternalBodyAnchorsKo(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		in := &ast.Internal{}
		children := []scan.Item{
			{Text: "삭제됨", No: 1, Struck: true},
			{Text: "`물 붓기 가능` 진행한다", No: 2, Span: scan.Span{Off: 0, Len: 10}},
			{Text: "`온도 확인` 판정한다", No: 3, Span: scan.Span{Off: 10, Len: 10}},
		}
		anchors, err := internalBodyAnchorsKo(in, 0, children, "t.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(anchors) != 2 {
			t.Fatalf("anchors = %d, want 2", len(anchors))
		}
		if len(in.Runs) != 1 || in.Runs[0] != "물 붓기 가능" {
			t.Fatalf("runs = %v", in.Runs)
		}
		if len(in.Checks) != 1 || in.Checks[0] != "온도 확인" {
			t.Fatalf("checks = %v", in.Checks)
		}
	})

	t.Run("MalformedError", func(t *testing.T) {
		in := &ast.Internal{}
		children := []scan.Item{{Text: "진행한다", No: 1}}
		if _, err := internalBodyAnchorsKo(in, 0, children, "t.md"); err == nil {
			t.Fatal("expected error for a run clause missing its backtick name")
		}
	})

	t.Run("UnrecognizedError", func(t *testing.T) {
		in := &ast.Internal{}
		children := []scan.Item{{Text: "알 수 없는 문장", No: 1}}
		if _, err := internalBodyAnchorsKo(in, 0, children, "t.md"); err == nil {
			t.Fatal("expected error for an unrecognized internal-body statement")
		}
	})
}
