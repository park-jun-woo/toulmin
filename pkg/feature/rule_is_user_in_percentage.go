//ff:func feature=feature type=rule control=sequence
//ff:what IsUserInPercentage: backing(float64)으로 지정된 비율 내에 사용자가 포함되는지 판정
package feature

// IsUserInPercentage checks if the user falls within the rollout percentage.
// backing is float64 (e.g., 0.3 = 30%). Uses deterministic hash, not rand.
func IsUserInPercentage(claim any, ground any, backing any) (bool, any) {
	ctx := ground.(*UserContext)
	pct := backing.(float64)
	return hashPercentage(ctx.ID) < pct, nil
}
