//ff:type feature=moderate type=model
//ff:what Action: 모더레이션 판정 결과 액션 (allow/flag/block)
package moderate

// Action represents the moderation decision.
type Action string

const (
	ActionAllow Action = "allow" // verdict > 0.3
	ActionFlag  Action = "flag"  // 0 < verdict <= 0.3
	ActionBlock Action = "block" // verdict <= 0
)
