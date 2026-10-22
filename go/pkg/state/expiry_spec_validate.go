//ff:func feature=state type=model control=sequence
//ff:what ExpirySpec.Validate: 필수 필드 검증
package state

// Validate always returns nil — ExpirySpec validation is context-dependent.
func (b *ExpirySpec) Validate() error { return nil }
