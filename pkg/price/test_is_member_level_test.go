//ff:func feature=price type=rule control=iteration dimension=1
//ff:what TestIsMemberLevel — tests IsMemberLevel rule
package price

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestIsMemberLevel(t *testing.T) {
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
			ctx := toulmin.NewContext()
			ctx.Set("membership", tt.membership)
			mb := &MemberSpec{Level: tt.level, Discount: &DiscountSpec{Name: tt.level, Rate: 0.1}}
			got, _ := IsMemberLevel(ctx, toulmin.Specs{mb})
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
