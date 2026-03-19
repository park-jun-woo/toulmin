//ff:func feature=price type=rule control=sequence
//ff:what IsAlreadyDiscounted: 이미 할인이 적용되었는지 판정
package price

// IsAlreadyDiscounted checks if the purchase is already discounted.
func IsAlreadyDiscounted(claim any, ground any, backing any) (bool, any) {
	req := claim.(*PurchaseRequest)
	discounted, _ := req.Metadata["discounted"].(bool)
	return discounted, nil
}
