//ff:func feature=tangl type=engine control=sequence
//ff:what TestCompileOperandValue — tests Path/SpecRef/Call operand compilation and errors
package eval

import (
	"errors"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestCompileOperandValue(t *testing.T) {
	boom := errors.New("boom")
	res := fakeResolver{values: map[string]func(toulmin.Context, toulmin.Specs) (any, error){
		"\x00ok":  func(toulmin.Context, toulmin.Specs) (any, error) { return 7, nil },
		"\x00bad": func(toulmin.Context, toulmin.Specs) (any, error) { return nil, boom },
	}}
	sc := newScope(testCtx(), toulmin.Specs{&testSpec{Name: "discount", MinPrice: 100}})

	t.Run("Path", func(t *testing.T) {
		vf, err := compileOperandValue(ast.Path{Steps: []string{"x"}}, res)
		if err != nil {
			t.Fatalf("compile: %v", err)
		}
		if v, ok, verr := vf(sc); !ok || verr != nil || v != "1" {
			t.Fatalf("path value: got %v,%v,%v", v, ok, verr)
		}
	})
	t.Run("SpecRef", func(t *testing.T) {
		vf, err := compileOperandValue(ast.SpecRef{Steps: []string{"discount", "minPrice"}}, res)
		if err != nil {
			t.Fatalf("compile: %v", err)
		}
		if v, ok, verr := vf(sc); !ok || verr != nil || v != float64(100) {
			t.Fatalf("specref value: got %v,%v,%v", v, ok, verr)
		}
	})
	t.Run("CallResolved", func(t *testing.T) {
		vf, err := compileOperandValue(ast.Call{Fn: ast.Ref{Name: "ok"}}, res)
		if err != nil {
			t.Fatalf("compile: %v", err)
		}
		if v, ok, verr := vf(sc); !ok || verr != nil || v != 7 {
			t.Fatalf("call value: got %v,%v,%v", v, ok, verr)
		}
	})
	t.Run("CallRuntimeError", func(t *testing.T) {
		vf, err := compileOperandValue(ast.Call{Fn: ast.Ref{Name: "bad"}}, res)
		if err != nil {
			t.Fatalf("compile: %v", err)
		}
		if _, ok, verr := vf(sc); ok || verr != boom {
			t.Fatalf("erroring call: got ok=%v err=%v", ok, verr)
		}
	})
	t.Run("CallUnresolved", func(t *testing.T) {
		if _, err := compileOperandValue(ast.Call{Fn: ast.Ref{Name: "nope"}}, res); err == nil {
			t.Fatal("want unresolved value error")
		}
	})
	t.Run("UnknownOperand", func(t *testing.T) {
		if _, err := compileOperandValue(nil, res); err == nil {
			t.Fatal("want unsupported-operand error")
		}
	})
}
