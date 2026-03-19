package price

import "testing"

func TestHasCoupon(t *testing.T) {
	db := &DiscountBacking{Name: "SAVE30", Rate: 0.3}
	tests := []struct {
		name    string
		base    float64
		coupons []Coupon
		want    bool
	}{
		{"has coupon, meets min", 10000, []Coupon{{Code: "SAVE30", MinPrice: 5000}}, true},
		{"has coupon, under min", 3000, []Coupon{{Code: "SAVE30", MinPrice: 5000}}, false},
		{"no coupons", 10000, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &PurchaseRequest{BasePrice: tt.base}
			ctx := &PriceContext{Coupons: tt.coupons}
			got, _ := HasCoupon(req, ctx, db)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsMemberLevel(t *testing.T) {
	tests := []struct {
		name       string
		membership string
		backing    string
		want       bool
	}{
		{"basic match", "basic", "basic", true},
		{"vip match", "vip", "vip", true},
		{"mismatch", "basic", "gold", false},
		{"none", "none", "basic", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &PriceContext{User: &User{Membership: tt.membership}}
			db := &DiscountBacking{Name: tt.backing, Rate: 0.1}
			got, _ := IsMemberLevel(nil, ctx, db)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

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

func TestIsAlreadyDiscounted(t *testing.T) {
	tests := []struct {
		name string
		meta map[string]any
		want bool
	}{
		{"discounted", map[string]any{"discounted": true}, true},
		{"not discounted", map[string]any{}, false},
		{"nil meta", nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &PurchaseRequest{Metadata: tt.meta}
			got, _ := IsAlreadyDiscounted(req, nil, nil)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsBulkOrder(t *testing.T) {
	tests := []struct {
		name string
		qty  int
		min  int
		want bool
	}{
		{"bulk", 100, 50, true},
		{"equal", 50, 50, true},
		{"not bulk", 10, 50, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &PurchaseRequest{Quantity: tt.qty}
			got, _ := IsBulkOrder(req, nil, tt.min)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
