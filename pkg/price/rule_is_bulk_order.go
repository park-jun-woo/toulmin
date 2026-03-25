//ff:func feature=price type=rule control=sequence
//ff:what IsBulkOrder: backing(int)으로 지정된 최소 수량 이상인지 판정
package price

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// IsBulkOrder checks if the order quantity meets the minimum specified by backing.
func IsBulkOrder(ctx toulmin.Context, backing toulmin.Backing) (bool, any) {
	quantity, _ := ctx.Get("quantity")
	bb := backing.(*BulkOrderBacking)
	return quantity.(int) >= bb.MinQuantity, nil
}
