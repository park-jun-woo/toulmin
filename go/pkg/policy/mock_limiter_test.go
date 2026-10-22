//ff:type feature=policy type=model
//ff:what mockLimiter — test helper rate limiter mock
package policy

type mockLimiter struct {
	limited map[string]bool
}
