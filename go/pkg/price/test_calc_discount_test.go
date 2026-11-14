//ff:func feature=price type=engine control=iteration dimension=1
//ff:what TestCalcDiscount — covers Min/Max clamp branches of calcDiscount for all short-circuit combinations
package price

import "testing"

func TestCalcDiscount(t *testing.T) {
	tests := []struct {
		name      string
		basePrice float64
		spec      *DiscountSpec
		want      float64
	}{
		// Min <= 0 (skip clamp), Max <= 0 (skip clamp): plain rate+fixed.
		{"no min no max", 100, &DiscountSpec{Rate: 0.1, Fixed: 5}, 15},
		// Min > 0 but d >= Min: clamp not applied.
		{"min set but not triggered", 100, &DiscountSpec{Rate: 0.5, Fixed: 0, Min: 10}, 50},
		// Min > 0 and d < Min: clamp applied.
		{"min set and triggered", 100, &DiscountSpec{Rate: 0.01, Fixed: 0, Min: 10}, 10},
		// Max > 0 but d <= Max: clamp not applied.
		{"max set but not triggered", 100, &DiscountSpec{Rate: 0.1, Fixed: 0, Max: 50}, 10},
		// Max > 0 and d > Max: clamp applied.
		{"max set and triggered", 100, &DiscountSpec{Rate: 0.9, Fixed: 0, Max: 20}, 20},
		// Both Min and Max set, neither triggered.
		{"min and max set neither triggered", 100, &DiscountSpec{Rate: 0.2, Fixed: 0, Min: 5, Max: 50}, 20},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calcDiscount(tt.basePrice, tt.spec)
			if got != tt.want {
				t.Errorf("calcDiscount(%v, %+v) = %v, want %v", tt.basePrice, tt.spec, got, tt.want)
			}
		})
	}
}
