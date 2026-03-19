//ff:type feature=approve type=model
//ff:what ApprovalRequest: 승인 요청
package approve

// ApprovalRequest represents an approval request.
type ApprovalRequest struct {
	Amount      float64
	Category    string
	RequesterID string
	Description string
	Metadata    map[string]any
}
