//ff:func feature=price type=model control=sequence
//ff:what BulkOrderSpec.SpecName: spec 타입 식별자 반환
package price

// SpecName returns the type identifier for BulkOrderSpec.
func (b *BulkOrderSpec) SpecName() string {
	return "BulkOrderSpec"
}
