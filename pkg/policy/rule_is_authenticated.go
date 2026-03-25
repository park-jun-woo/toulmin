//ff:func feature=policy type=rule control=sequence
//ff:what IsAuthenticated: 사용자가 인증되었는지 판정
package policy

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// IsAuthenticated returns true if the request has a non-nil user.
func IsAuthenticated(ctx toulmin.Context, specs toulmin.Specs) (bool, any) {
	user, _ := ctx.Get("user")
	return user != nil, nil
}
