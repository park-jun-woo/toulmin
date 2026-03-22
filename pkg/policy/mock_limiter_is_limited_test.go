//ff:func feature=policy type=model control=sequence
//ff:what mockLimiter.IsLimited — checks if key is rate limited
package policy

func (m *mockLimiter) IsLimited(key string) bool { return m.limited[key] }
