//ff:func feature=moderate type=model control=sequence
//ff:what ClassifierSpec.SpecName: spec 타입 식별자 반환
package moderate

// SpecName returns the type identifier for ClassifierSpec.
func (b *ClassifierSpec) SpecName() string {
	return "ClassifierSpec"
}
