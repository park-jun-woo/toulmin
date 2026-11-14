//ff:func feature=state type=model control=sequence
//ff:what OwnerSpec.Validate: 필수 필드 검증
package state

// Validate always returns nil — OwnerSpec has no required fields.
func (b *OwnerSpec) Validate() error { return nil }
