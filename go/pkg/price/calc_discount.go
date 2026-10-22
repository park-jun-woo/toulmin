//ff:func feature=price type=engine control=sequence
//ff:what calcDiscount: 단일 DiscountSpec의 할인액 계산
package price

// calcDiscount computes the discount for a single DiscountSpec.
func calcDiscount(basePrice float64, db *DiscountSpec) float64 {
	d := basePrice*db.Rate + db.Fixed
	if db.Min > 0 && d < db.Min {
		d = db.Min
	}
	if db.Max > 0 && d > db.Max {
		d = db.Max
	}
	return d
}
