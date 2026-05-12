//ff:func feature=feature type=rule control=sequence
//ff:what IsBetaUser: 베타 사용자인지 판정
package feature

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// IsBetaUser returns true if the user has the "beta" attribute set to true.
func IsBetaUser(ctx toulmin.Context, specs toulmin.Specs) (bool, any) {
	attributes, _ := ctx.Get("attributes")
	attrs, ok := attributes.(map[string]any)
	if !ok {
		return false, nil
	}
	beta, _ := attrs["beta"].(bool)
	return beta, nil
}
