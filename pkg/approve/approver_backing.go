//ff:type feature=approve type=model
//ff:what ApproverBacking: 결재자 관련 rule의 backing 타입
package approve

// ApproverBacking carries extraction functions for approver checks.
type ApproverBacking struct {
	Role     string              // role to match (for HasApprovalRole)
	Level    int                 // minimum level (for IsAboveLevel)
	IDFunc   func(any) string   // extracts approver ID
	RoleFunc func(any) string   // extracts approver role
	LevelFunc func(any) int    // extracts approver level
}
