//ff:func feature=price type=rule control=sequence
//ff:what IsMemberLevel: backing(MemberBacking)으로 사용자 멤버십 등급 비교
package price

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// IsMemberLevel checks if the user's membership matches backing.Level.
// Returns the Discount as evidence if matched.
func IsMemberLevel(ctx toulmin.Context, backing toulmin.Backing) (bool, any) {
	membership, _ := ctx.Get("membership")
	mb := backing.(*MemberBacking)
	ms := membership.(string)
	if ms == "" {
		return false, nil
	}
	if ms == mb.Level {
		return true, mb.Discount
	}
	return false, nil
}
