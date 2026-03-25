//ff:func feature=price type=model control=sequence
//ff:what DiscountSpec.SpecName: spec 타입 식별자 반환
package price

// SpecName returns the type identifier for DiscountSpec.
func (b *DiscountSpec) SpecName() string {
	return "DiscountSpec"
}
