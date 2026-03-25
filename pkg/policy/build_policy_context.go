//ff:func feature=policy type=adapter control=sequence
//ff:what buildPolicyContext: RequestContext → toulmin.Context 변환
package policy

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// buildPolicyContext converts a RequestContext into a toulmin.Context.
func buildPolicyContext(rc *RequestContext) toulmin.Context {
	ctx := toulmin.NewContext()
	ctx.Set("user", rc.User)
	ctx.Set("clientIP", rc.ClientIP)
	ctx.Set("resourceOwnerID", rc.ResourceOwnerID)
	ctx.Set("headers", rc.Headers)
	ctx.Set("rateLimiter", rc.RateLimiter)
	ctx.Set("metadata", rc.Metadata)
	ctx.Set("role", rc.Role)
	ctx.Set("userID", rc.UserID)
	ctx.Set("resourceOwner", rc.ResourceOwner)
	ctx.Set("ipBlocked", rc.IPBlocked)
	return ctx
}
