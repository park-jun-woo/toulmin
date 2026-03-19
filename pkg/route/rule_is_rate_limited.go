//ff:func feature=route type=rule control=sequence
//ff:what IsRateLimited: backing(RateLimiter)으로 전달된 limiter로 rate limit 판정
package route

// IsRateLimited checks if the client is rate limited.
// backing is a RateLimiter instance.
func IsRateLimited(claim any, ground any, backing any) (bool, any) {
	ctx := ground.(*RouteContext)
	limiter := backing.(RateLimiter)
	return limiter.IsLimited(ctx.ClientIP), nil
}
