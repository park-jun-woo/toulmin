//ff:func feature=approve type=rule control=sequence
//ff:what IsUrgent: 긴급 승인 요청인지 판정
package approve

// IsUrgent checks if the request is marked as urgent.
func IsUrgent(claim any, ground any, backing any) (bool, any) {
	req := claim.(*ApprovalRequest)
	urgent, _ := req.Metadata["urgent"].(bool)
	return urgent, nil
}
