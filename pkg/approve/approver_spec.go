//ff:type feature=approve type=model
//ff:what ApproverSpec: 결재자 관련 rule의 spec 타입
package approve

// ApproverSpec carries criteria for approver checks.
type ApproverSpec struct {
	Role  string // role to match (for HasApprovalRole)
	Level int    // minimum level (for IsAboveLevel)
}
