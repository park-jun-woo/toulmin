//ff:func feature=tangl type=engine control=sequence
//ff:what TestSpecVal — tests SpecRef resolution to the Spec itself and into a Spec field
package eval

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestSpecVal(t *testing.T) {
	sp := &testSpec{Name: "discount", MinPrice: 100}
	sc := newScope(toulmin.NewContext(), toulmin.Specs{sp})

	t.Run("SpecItself", func(t *testing.T) {
		v, ok := specVal(ast.SpecRef{Steps: []string{"discount"}}, sc)
		if !ok || v != toulmin.Spec(sp) {
			t.Fatalf("want the spec,true got %v,%v", v, ok)
		}
	})
	t.Run("SpecField", func(t *testing.T) {
		v, ok := specVal(ast.SpecRef{Steps: []string{"discount", "minPrice"}}, sc)
		if !ok || v != float64(100) {
			t.Fatalf("want 100,true got %v,%v", v, ok)
		}
	})
	t.Run("EmptySteps", func(t *testing.T) {
		if _, ok := specVal(ast.SpecRef{}, sc); ok {
			t.Fatal("expected ok=false for empty steps")
		}
	})
	t.Run("UnknownSpec", func(t *testing.T) {
		if _, ok := specVal(ast.SpecRef{Steps: []string{"missing"}}, sc); ok {
			t.Fatal("expected ok=false for unknown spec")
		}
	})
}
