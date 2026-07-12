//ff:func feature=tangl type=codegen control=sequence
//ff:what equivDiscount.SpecName — returns the "discount" spec identifier
package gen

// SpecName returns the spec term name the cross-test document references.
func (s *equivDiscount) SpecName() string { return "discount" }
