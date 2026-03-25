//ff:func feature=feature type=rule control=sequence
//ff:what IsBetaUser: 베타 사용자인지 판정
package feature

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// IsBetaUser returns true if the user has the "beta" attribute set to true.
func IsBetaUser(ctx toulmin.Context, specs toulmin.Specs) (bool, any) {
	attributes, _ := ctx.Get("attributes")
	beta, _ := attributes.(map[string]any)["beta"].(bool)
	return beta, nil
}
