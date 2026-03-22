//ff:func feature=price type=rule control=iteration dimension=1
//ff:what TestHasActivePromotion — tests HasActivePromotion rule
package price

import "testing"

func TestHasActivePromotion(t *testing.T) {
	db := &DiscountBacking{Name: "blackfriday", Fixed: 5000}
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
			ctx := &PriceContext{Promotions: tt.promos}
			got, _ := HasActivePromotion(nil, ctx, db)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
