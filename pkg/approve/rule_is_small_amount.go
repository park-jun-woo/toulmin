//ff:func feature=approve type=rule control=sequence
//ff:what IsSmallAmount: backing(float64)으로 지정된 금액 임계값 이하인지 판정
package approve

// IsSmallAmount checks if the requested amount is at or below backing (float64).
func IsSmallAmount(claim any, ground any, backing any) (bool, any) {
	req := claim.(*ApprovalRequest)
	threshold := backing.(float64)
	return req.Amount <= threshold, nil
}
