//ff:func feature=price type=rule control=sequence
//ff:what TestIsMemberLevel_NilUser — tests IsMemberLevel with nil user
package price

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestIsMemberLevel_NilUser(t *testing.T) {
	ctx := toulmin.NewContext()
	ctx.Set("membership", "")
	mb := &MemberSpec{Level: "basic", Discount: &DiscountSpec{}}
	got, _ := IsMemberLevel(ctx, toulmin.Specs{mb})
	if got {
		t.Error("expected false for nil user")
	}
}
