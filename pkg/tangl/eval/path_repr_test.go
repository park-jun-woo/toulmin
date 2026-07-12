//ff:func feature=tangl type=engine control=sequence
//ff:what TestPathRepr — tests dotted rendering and binder-indexed base/rest forms
package eval

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestPathRepr(t *testing.T) {
	sc := newScope(toulmin.NewContext(), nil)

	if got := pathRepr(ast.Path{}, sc); got != "" {
		t.Fatalf("empty path: got %q", got)
	}
	if got := pathRepr(ast.Path{Steps: []string{"author", "trustScore"}}, sc); got != "author.trustScore" {
		t.Fatalf("dotted path: got %q", got)
	}
	sc.binders["c"] = &binding{source: ast.Path{Steps: []string{"coupons"}}, index: 1}
	if got := pathRepr(ast.Path{Steps: []string{"c"}}, sc); got != "coupons[1]" {
		t.Fatalf("bare binder: got %q", got)
	}
	if got := pathRepr(ast.Path{Steps: []string{"c", "minPrice"}}, sc); got != "coupons[1].minPrice" {
		t.Fatalf("binder with rest: got %q", got)
	}
}
