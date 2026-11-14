//ff:func feature=tangl type=parser control=sequence
//ff:what TestEndpointBodyAnchorsKo — struck skip, require/run/check dispatch, and each error path
package ko

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/codec/scan"
)

func TestEndpointBodyAnchorsKo(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		ep := &ast.Endpoint{}
		children := []scan.Item{
			{Text: "삭제됨", No: 1, Struck: true},
			{Text: "`사용자` 는 필수다", No: 2, Span: scan.Span{Off: 0, Len: 10}},
			{Text: "`물 붓기 가능` 진행한다", No: 3, Span: scan.Span{Off: 10, Len: 10}},
			{Text: "`온도 확인` 판정한다", No: 4, Span: scan.Span{Off: 20, Len: 10}},
		}
		anchors, err := endpointBodyAnchorsKo(ep, 0, children, "t.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(anchors) != 3 {
			t.Fatalf("anchors = %d, want 3", len(anchors))
		}
		if len(ep.Requires) != 1 || len(ep.Runs) != 1 || len(ep.Checks) != 1 {
			t.Fatalf("req=%d run=%d check=%d", len(ep.Requires), len(ep.Runs), len(ep.Checks))
		}
	})

	t.Run("RequireError", func(t *testing.T) {
		ep := &ast.Endpoint{}
		children := []scan.Item{{Text: "`f` 는 필수다, 이상한말", No: 1}}
		if _, err := endpointBodyAnchorsKo(ep, 0, children, "t.md"); err == nil {
			t.Fatal("expected error for a malformed 필수다 tail")
		}
	})

	t.Run("RunCheckError", func(t *testing.T) {
		ep := &ast.Endpoint{}
		children := []scan.Item{{Text: "진행한다", No: 1}}
		if _, err := endpointBodyAnchorsKo(ep, 0, children, "t.md"); err == nil {
			t.Fatal("expected error for a run clause missing its backtick name")
		}
	})

	t.Run("UnrecognizedError", func(t *testing.T) {
		ep := &ast.Endpoint{}
		children := []scan.Item{{Text: "알 수 없는 문장", No: 1}}
		if _, err := endpointBodyAnchorsKo(ep, 0, children, "t.md"); err == nil {
			t.Fatal("expected error for an unrecognized 업무 statement")
		}
	})
}
