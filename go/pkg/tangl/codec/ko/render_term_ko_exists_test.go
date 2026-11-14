//ff:func feature=tangl type=parser control=sequence
//ff:what TestRenderTermKoExists — a nested-body exists renders a head plus indented terms; source and body errors propagate
package ko

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestRenderTermKoExists(t *testing.T) {
	cmp := ast.Compare{Lhs: ast.Path{Steps: []string{"파일", "이름"}}, Op: "is not empty"}
	ex := ast.Exists{Source: ast.Path{Steps: []string{"파일들"}}, Binder: "파일", Cond: ast.Logic{Op: "and", Terms: []ast.Expr{cmp, cmp}}}
	lines, err := renderTermKo("", ex, 1)
	if err != nil || len(lines) != 3 {
		t.Fatalf("lines=%q err=%v", lines, err)
	}
	if !strings.HasPrefix(lines[0], "  - `파일들` 중 `파일` 이 다음인 것이 있으면") || !strings.HasPrefix(lines[1], "    - ") {
		t.Fatalf("lines = %q", lines)
	}

	noSource := ex
	noSource.Source = nil
	if _, err := renderTermKo("", noSource, 0); err == nil || !strings.Contains(err.Error(), "unsupported operand") {
		t.Fatalf("nil source: %v", err)
	}
	badBody := ex
	badBody.Cond = ast.Logic{Op: "and", Terms: []ast.Expr{cmp, nil}}
	if _, err := renderTermKo("", badBody, 0); err == nil || !strings.Contains(err.Error(), "unsupported condition term") {
		t.Fatalf("bad body: %v", err)
	}
}
