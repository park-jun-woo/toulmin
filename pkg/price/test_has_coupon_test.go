//ff:func feature=price type=rule control=iteration dimension=1
//ff:what TestHasCoupon — tests HasCoupon rule
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
