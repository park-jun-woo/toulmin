//ff:type feature=approve type=model
//ff:what Approver: 결재자 정보
package approve

// Approver represents an approver's identity and level.
type Approver struct {
	ID    string
	Role  string
	Level int
}
