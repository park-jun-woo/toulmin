//ff:func feature=price type=rule control=sequence
//ff:what IsBulkOrder: backing(int)으로 지정된 최소 수량 이상인지 판정
package price

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// IsBulkOrder checks if the order quantity meets the minimum specified by backing.
func IsBulkOrder(claim any, ground any, backing toulmin.Backing) (bool, any) {
	req := claim.(*PurchaseRequest)
	bb := backing.(*BulkOrderBacking)
	return req.Quantity >= bb.MinQuantity, nil
}
