//ff:type feature=price type=model
//ff:what BulkOrderBacking: IsBulkOrder rule의 backing 타입
package price

// BulkOrderBacking carries minimum quantity criteria for bulk order checks.
type BulkOrderBacking struct {
	MinQuantity int // minimum quantity for bulk order
}
