//ff:func feature=tangl type=engine control=sequence
//ff:what TestCompileHostRule — host Rules declarations bind through Resolver and preserve evidence
package eval

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/go/pkg/toulmin"
)

func TestCompileHostRule(t *testing.T) {
	res := fakeResolver{predicates: map[string]func(toulmin.Context, toulmin.Specs) (bool, any){
		"check\x00bad": func(toulmin.Context, toulmin.Specs) (bool, any) { return true, "host evidence" },
	}}
	fn, err := Compile(ast.InlineRule{Name: "bad", Using: &ast.Ref{Alias: "check", Name: "bad"}}, res)
	if err != nil {
		t.Fatal(err)
	}
	got, evidence := fn(toulmin.NewContext(), nil)
	if !got || evidence != "host evidence" {
		t.Fatalf("got=%v evidence=%v", got, evidence)
	}
}
