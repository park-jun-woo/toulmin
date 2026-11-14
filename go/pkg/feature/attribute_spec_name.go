//ff:func feature=feature type=model control=sequence
//ff:what AttributeSpec.SpecName: spec 타입 식별자 반환
package feature

// SpecName returns the type identifier for AttributeSpec.
func (b *AttributeSpec) SpecName() string {
	return "AttributeSpec"
}
