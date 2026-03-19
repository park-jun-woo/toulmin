//ff:func feature=price type=rule control=sequence
//ff:what IsBulkOrder: backing(int)으로 지정된 최소 수량 이상인지 판정
package price

// IsBulkOrder checks if the order quantity meets the minimum specified by backing (int).
func IsBulkOrder(claim any, ground any, backing any) (bool, any) {
	req := claim.(*PurchaseRequest)
	minQty := backing.(int)
	return req.Quantity >= minQty, nil
}
