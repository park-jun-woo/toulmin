//ff:func feature=state type=model control=sequence
//ff:what ExpiryBacking.Validate: 필수 필드 검증
package state

// Validate always returns nil — ExpiryBacking validation is context-dependent.
func (b *ExpiryBacking) Validate() error { return nil }
