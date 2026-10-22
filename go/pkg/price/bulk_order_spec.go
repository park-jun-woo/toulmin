//ff:type feature=price type=model
//ff:what BulkOrderSpec: IsBulkOrder rule의 spec 타입
package price

// BulkOrderSpec carries minimum quantity criteria for bulk order checks.
type BulkOrderSpec struct {
	MinQuantity int // minimum quantity for bulk order
}
