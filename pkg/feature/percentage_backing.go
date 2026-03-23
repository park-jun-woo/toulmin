//ff:type feature=feature type=model
//ff:what PercentageBacking: IsUserInPercentage rule의 backing 타입
package feature

// PercentageBacking carries rollout percentage criteria for feature flag checks.
type PercentageBacking struct {
	Percentage float64 // rollout percentage (0.0 ~ 1.0)
}
