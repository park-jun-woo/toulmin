//ff:func feature=price type=model control=sequence
//ff:what DiscountBacking.BackingName: backing 타입 식별자 반환
package price

// BackingName returns the type identifier for DiscountBacking.
func (b *DiscountBacking) BackingName() string {
	return "DiscountBacking"
}
