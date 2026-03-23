//ff:func feature=state type=model control=sequence
//ff:what OwnerBacking.Validate: 필수 필드 검증
package state

// Validate always returns nil — OwnerBacking has no required fields.
func (b *OwnerBacking) Validate() error { return nil }
