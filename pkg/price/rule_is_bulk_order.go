//ff:func feature=price type=rule control=sequence
//ff:what IsBulkOrder: spec(int)으로 지정된 최소 수량 이상인지 판정
package price

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// IsBulkOrder checks if the order quantity meets the minimum specified by spec.
func IsBulkOrder(ctx toulmin.Context, specs toulmin.Specs) (bool, any) {
	quantity, _ := ctx.Get("quantity")
	if len(specs) == 0 {
		return false, nil
	}
	bb := specs[0].(*BulkOrderSpec)
	return quantity.(int) >= bb.MinQuantity, nil
}
