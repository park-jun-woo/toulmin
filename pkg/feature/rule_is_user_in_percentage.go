//ff:func feature=feature type=rule control=sequence
//ff:what IsUserInPercentage: spec(PercentageSpec)으로 지정된 비율 내에 사용자가 포함되는지 판정
package feature

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// IsUserInPercentage checks if the user falls within the rollout percentage.
// Uses deterministic hash, not rand.
func IsUserInPercentage(ctx toulmin.Context, specs toulmin.Specs) (bool, any) {
	id, _ := ctx.Get("id")
	pb := specs[0].(*PercentageSpec)
	return hashPercentage(id.(string)) < pb.Percentage, nil
}
