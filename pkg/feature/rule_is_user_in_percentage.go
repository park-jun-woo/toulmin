//ff:func feature=feature type=rule control=sequence
//ff:what IsUserInPercentage: backing(PercentageBacking)으로 지정된 비율 내에 사용자가 포함되는지 판정
package feature

// IsUserInPercentage checks if the user falls within the rollout percentage.
// Uses deterministic hash, not rand.
func IsUserInPercentage(claim any, ground any, backing any) (bool, any) {
	ctx := ground.(*UserContext)
	pb := backing.(*PercentageBacking)
	return hashPercentage(ctx.ID) < pb.Percentage, nil
}
