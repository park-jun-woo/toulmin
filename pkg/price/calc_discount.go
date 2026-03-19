//ff:func feature=price type=engine control=sequence
//ff:what calcDiscount: 단일 DiscountBacking의 할인액 계산
package price

// calcDiscount computes the discount for a single DiscountBacking.
func calcDiscount(basePrice float64, db *DiscountBacking) float64 {
	d := basePrice*db.Rate + db.Fixed
	if db.Min > 0 && d < db.Min {
		d = db.Min
	}
	if db.Max > 0 && d > db.Max {
		d = db.Max
	}
	return d
}
