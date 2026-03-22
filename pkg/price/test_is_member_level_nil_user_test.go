//ff:func feature=price type=rule control=sequence
//ff:what TestIsMemberLevel_NilUser — tests IsMemberLevel with nil user
package price

import "testing"

func TestIsMemberLevel_NilUser(t *testing.T) {
	memberFunc := func(u any) string { return u.(*testUser).Membership }
	ctx := &PriceContext{User: nil}
	mb := &MemberBacking{Level: "basic", MembershipFunc: memberFunc, Discount: &DiscountBacking{}}
	got, _ := IsMemberLevel(nil, ctx, mb)
	if got {
		t.Error("expected false for nil user")
	}
}
