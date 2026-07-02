//ff:func feature=tangl type=codegen control=sequence
//ff:what TestRequireFields — tests requireFields for empty and multi-item input branches
package gen

import (
	"reflect"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestRequireFields(t *testing.T) {
	t.Run("empty input", func(t *testing.T) {
		got := requireFields(nil)
		if len(got) != 0 {
			t.Errorf("expected empty slice, got %v", got)
		}
	})

	t.Run("multiple requires extracted in order", func(t *testing.T) {
		reqs := []ast.Require{
			{Field: "amount"},
			{Field: "label"},
		}
		got := requireFields(reqs)
		want := []string{"amount", "label"}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}
