//ff:func feature=price type=rule control=iteration dimension=1
//ff:what TestIsMemberLevel — tests IsMemberLevel rule
package price

import "testing"

func TestIsMemberLevel(t *testing.T) {
	memberFunc := func(u any) string { return u.(*testUser).Membership }
	tests := []struct {
		name       string
		membership string
		level      string
		want       bool
	}{
		{"basic match", "basic", "basic", true},
		{"vip match", "vip", "vip", true},
		{"mismatch", "basic", "gold", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &PriceContext{User: &testUser{Membership: tt.membership}}
			mb := &MemberBacking{Level: tt.level, MembershipFunc: memberFunc, Discount: &DiscountBacking{Name: tt.level, Rate: 0.1}}
			got, _ := IsMemberLevel(nil, ctx, mb)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
