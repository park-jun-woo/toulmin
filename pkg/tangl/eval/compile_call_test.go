//ff:func feature=tangl type=engine control=sequence
//ff:what TestCompileCall — tests predicate host-evidence passthrough, synthesized fact, and value-call error
package eval

import (
	"errors"
	"reflect"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestCompileCall(t *testing.T) {
	t.Run("PredicateHostEvidencePassThrough", func(t *testing.T) {
		sentinel := toulmin.Fact{Where: "host", Expected: "custom", Actual: 7}
		res := fakeResolver{predicates: map[string]func(toulmin.Context, toulmin.Specs) (bool, any){
			"\x00spam": func(toulmin.Context, toulmin.Specs) (bool, any) { return true, sentinel },
		}}
		b, ev := runRule(t, ast.Call{Fn: ast.Ref{Name: "spam"}}, res, testCtx())
		if !b || !reflect.DeepEqual(ev, sentinel) {
			t.Fatalf("want host evidence passthrough, got %v %+v", b, ev)
		}
	})
	t.Run("PredicateNilEvidenceSynthesized", func(t *testing.T) {
		res := fakeResolver{predicates: map[string]func(toulmin.Context, toulmin.Specs) (bool, any){
			"\x00spam": func(toulmin.Context, toulmin.Specs) (bool, any) { return true, nil },
		}}
		b, ev := runRule(t, ast.Call{Fn: ast.Ref{Name: "spam"}}, res, testCtx())
		f := asFact(t, ev)
		if !b || f.Where != "spam" || f.Expected != "holds" || f.Actual != true {
			t.Fatalf("synthesized call fact mismatch: %v %+v", b, f)
		}
	})
	t.Run("ValueCallErrorInActual", func(t *testing.T) {
		boom := errors.New("boom")
		res := fakeResolver{values: map[string]func(toulmin.Context, toulmin.Specs) (any, error){
			"\x00v": func(toulmin.Context, toulmin.Specs) (any, error) { return nil, boom },
		}}
		cond := ast.Compare{Lhs: ast.Call{Fn: ast.Ref{Name: "v"}}, Op: "equals", Value: "x"}
		b, ev := runRule(t, cond, res, testCtx())
		f := asFact(t, ev)
		if b || f.Where != "v" || f.Expected != "equals x" || f.Actual != error(boom) {
			t.Fatalf("value-call-error fact mismatch: %v %+v", b, f)
		}
	})
	t.Run("UnresolvedBarePredicateIsCompileError", func(t *testing.T) {
		if _, err := Compile(ast.InlineRule{Cond: ast.Call{Fn: ast.Ref{Name: "ghost"}}}, fakeResolver{}); err == nil {
			t.Fatal("expected compile error for unresolved predicate")
		}
	})
}
