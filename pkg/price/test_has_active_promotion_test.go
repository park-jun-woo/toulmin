//ff:func feature=price type=rule control=iteration dimension=1
//ff:what TestHasActivePromotion — tests HasActivePromotion rule
package price

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestHasActivePromotion(t *testing.T) {
	db := &DiscountSpec{Name: "blackfriday", Fixed: 5000}
	tests := []struct {
		name   string
		promos []Promotion
		want   bool
	}{
		{"active", []Promotion{{Name: "blackfriday", Active: true}}, true},
		{"inactive", []Promotion{{Name: "blackfriday", Active: false}}, false},
		{"wrong name", []Promotion{{Name: "summer", Active: true}}, false},
		{"no promos", nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := toulmin.NewContext()
			ctx.Set("promotions", tt.promos)
			got, _ := HasActivePromotion(ctx, toulmin.Specs{db})
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
