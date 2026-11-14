//ff:func feature=tangl type=codegen control=sequence
//ff:what TestCallFnSelector — tests aliased verbatim vs bare goIdent call symbols
package gen

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestCallFnSelector(t *testing.T) {
	if got := callFnSelector(ast.Call{Fn: ast.Ref{Alias: "nlp", Name: "score"}}); got != "nlp.score" {
		t.Errorf("aliased call: got %q", got)
	}
	if got := callFnSelector(ast.Call{Fn: ast.Ref{Name: "is spam"}}); got != "isSpam" {
		t.Errorf("bare call should become a goIdent: got %q", got)
	}
}
