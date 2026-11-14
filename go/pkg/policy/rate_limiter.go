//ff:type feature=policy type=interface
//ff:what RateLimiter: rate limiting 추상화 인터페이스
package policy

// RateLimiter abstracts rate limiting logic.
type RateLimiter interface {
	IsLimited(key string) bool
}
