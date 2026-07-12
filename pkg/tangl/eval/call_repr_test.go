//ff:func feature=tangl type=util control=sequence
//ff:what TestCallRepr — tests Call rendering with and without an alias
package eval

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestCallRepr(t *testing.T) {
	if got := callRepr(ast.Call{Fn: ast.Ref{Alias: "nlp", Name: "score"}}); got != "nlp.score" {
		t.Fatalf("aliased call: got %q", got)
	}
	if got := callRepr(ast.Call{Fn: ast.Ref{Name: "score"}}); got != "score" {
		t.Fatalf("bare call: got %q", got)
	}
}
