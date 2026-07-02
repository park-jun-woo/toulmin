//ff:func feature=feature type=rule control=sequence
//ff:what TestHasAttribute_Branches — covers empty specs and non-map attributes branches of HasAttribute
package feature

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestHasAttribute_Branches(t *testing.T) {
	t.Run("EmptySpecs", func(t *testing.T) {
		ctx := toulmin.NewContext()
		ctx.Set("attributes", map[string]any{"plan": "pro"})
		got, _ := HasAttribute(ctx, toulmin.Specs{})
		if got != false {
			t.Errorf("got %v, want false", got)
		}
	})

	t.Run("NonMapAttributes", func(t *testing.T) {
		ctx := toulmin.NewContext()
		ctx.Set("attributes", "not-a-map")
		got, _ := HasAttribute(ctx, toulmin.Specs{&AttributeSpec{Key: "plan", Value: "pro"}})
		if got != false {
			t.Errorf("got %v, want false", got)
		}
	})
}
