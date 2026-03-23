//ff:func feature=price type=model control=sequence
//ff:what BulkOrderBacking.BackingName: backing 타입 식별자 반환
package price

// BackingName returns the type identifier for BulkOrderBacking.
func (b *BulkOrderBacking) BackingName() string {
	return "BulkOrderBacking"
}
