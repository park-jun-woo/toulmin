//ff:func feature=tangl type=codegen control=sequence
//ff:what equivDiscount.Validate — always valid cross-test spec
package gen

// Validate reports the cross-test spec as always valid.
func (s *equivDiscount) Validate() error { return nil }
