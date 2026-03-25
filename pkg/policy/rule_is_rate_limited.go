//ff:func feature=policy type=rule control=sequence
//ff:what IsRateLimited: 클라이언트가 rate limit에 걸렸는지 판정
package policy

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// IsRateLimited checks if the client is rate limited.
func IsRateLimited(claim any, ground any, backing toulmin.Backing) (bool, any) {
	ctx := ground.(*RequestContext)
	return ctx.RateLimiter.IsLimited(ctx.ClientIP), nil
}
