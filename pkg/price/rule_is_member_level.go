//ff:func feature=price type=rule control=sequence
//ff:what IsMemberLevel: spec(MemberSpec)으로 사용자 멤버십 등급 비교
package price

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// IsMemberLevel checks if the user's membership matches spec.Level.
// Returns the Discount as evidence if matched.
func IsMemberLevel(ctx toulmin.Context, specs toulmin.Specs) (bool, any) {
	membership, _ := ctx.Get("membership")
	if len(specs) == 0 {
		return false, nil
	}
	mb := specs[0].(*MemberSpec)
	ms := membership.(string)
	if ms == "" {
		return false, nil
	}
	if ms == mb.Level {
		return true, mb.Discount
	}
	return false, nil
}
