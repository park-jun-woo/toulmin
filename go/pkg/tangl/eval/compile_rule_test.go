//ff:func feature=tangl type=engine control=sequence
//ff:what TestCompileRule — expression rules compile; host rules need a resolver and a resolvable predicate
package eval

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/go/pkg/toulmin"
)

func TestCompileRule(t *testing.T) {
	expr := ast.InlineRule{Name: "big", Cond: ast.Compare{Field: "n", Op: "is greater than", Value: "3"}}
	fn, err := compileRule(expr, nil)
	if err != nil {
		t.Fatalf("expression rule: %v", err)
	}
	ctx := toulmin.NewContext()
	ctx.Set("n", 5)
	if got, _ := fn(newScope(ctx, nil)); !got {
		t.Fatal("expression rule: 5 > 3 should hold")
	}

	host := ast.InlineRule{Name: "spam", Using: &ast.Ref{Alias: "nlp", Name: "isSpam"}}
	if _, err := compileRule(host, nil); err == nil || !strings.Contains(err.Error(), "no resolver") {
		t.Fatalf("host rule without resolver: err = %v", err)
	}
	if _, err := compileRule(host, fakeResolver{}); err == nil || !strings.Contains(err.Error(), "unresolved host predicate nlp.isSpam") {
		t.Fatalf("unresolved host predicate: err = %v", err)
	}
	res := fakeResolver{predicates: map[string]func(toulmin.Context, toulmin.Specs) (bool, any){
		"nlp\x00isSpam": func(c toulmin.Context, s toulmin.Specs) (bool, any) { return len(s) == 1, "host" },
	}}
	fn, err = compileRule(host, res)
	if err != nil {
		t.Fatalf("resolved host rule: %v", err)
	}
	if got, ev := fn(newScope(ctx, toulmin.Specs{&testSpec{Name: "s"}})); !got || ev != "host" {
		t.Fatalf("host closure must forward ctx and specs: got %v,%v", got, ev)
	}
}
