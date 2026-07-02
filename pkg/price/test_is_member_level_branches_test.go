//ff:func feature=price type=rule control=sequence
//ff:what TestIsMemberLevel_Branches — covers empty specs and membership type-assertion failure branches of IsMemberLevel
package price

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestIsMemberLevel_Branches(t *testing.T) {
	t.Run("empty specs", func(t *testing.T) {
		ctx := toulmin.NewContext()
		ctx.Set("membership", "basic")

		got, evidence := IsMemberLevel(ctx, toulmin.Specs{})
		if got {
			t.Errorf("expected false for empty specs, got %v", got)
		}
		if evidence != nil {
			t.Errorf("expected nil evidence, got %v", evidence)
		}
	})

	t.Run("membership wrong type", func(t *testing.T) {
		ctx := toulmin.NewContext()
		ctx.Set("membership", 123)

		mb := &MemberSpec{Level: "basic", Discount: &DiscountSpec{}}
		got, _ := IsMemberLevel(ctx, toulmin.Specs{mb})
		if got {
			t.Errorf("expected false when membership is not a string, got %v", got)
		}
	})

	t.Run("membership unset", func(t *testing.T) {
		ctx := toulmin.NewContext()

		mb := &MemberSpec{Level: "basic", Discount: &DiscountSpec{}}
		got, _ := IsMemberLevel(ctx, toulmin.Specs{mb})
		if got {
			t.Errorf("expected false when membership is unset, got %v", got)
		}
	})
}
