//ff:func feature=feature type=rule control=sequence
//ff:what IsBetaUser: 베타 사용자인지 판정
package feature

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// IsBetaUser returns true if the user has the "beta" attribute set to true.
func IsBetaUser(claim any, ground any, backing toulmin.Backing) (bool, any) {
	ctx := ground.(*UserContext)
	beta, _ := ctx.Attributes["beta"].(bool)
	return beta, nil
}
