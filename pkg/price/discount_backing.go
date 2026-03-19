//ff:type feature=price type=model
//ff:what DiscountBacking: 할인 판정 기준 (Name, Rate, Fixed, Min, Max)
package price

// DiscountBacking carries discount criteria.
// Rate and Fixed can be combined. Min/Max constrain the result.
type DiscountBacking struct {
	Name  string  // discount name ("SAVE30", "basic", "blackfriday")
	Rate  float64 // percentage discount (0.3 = 30%). 0 = not applied
	Fixed float64 // fixed amount discount (5000). 0 = not applied
	Min   float64 // minimum discount guarantee. 0 = no minimum
	Max   float64 // maximum discount cap. 0 = no cap
}
