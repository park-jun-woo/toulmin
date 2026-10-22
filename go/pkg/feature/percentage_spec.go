//ff:type feature=feature type=model
//ff:what PercentageSpec: IsUserInPercentage rule의 spec 타입
package feature

// PercentageSpec carries rollout percentage criteria for feature flag checks.
type PercentageSpec struct {
	Percentage float64 // rollout percentage (0.0 ~ 1.0)
}
