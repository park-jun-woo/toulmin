//ff:func feature=moderate type=rule control=sequence
//ff:what IsVerifiedUser: 인증된 사용자인지 판정
package moderate

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// IsVerifiedUser returns true if the author is verified.
func IsVerifiedUser(ctx toulmin.Context, backing toulmin.Backing) (bool, any) {
	author, _ := ctx.Get("author")
	return author.(*Author).Verified, nil
}
