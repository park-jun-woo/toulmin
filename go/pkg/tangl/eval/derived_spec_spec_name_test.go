//ff:func feature=tangl type=engine control=sequence
//ff:what SpecName — derivedSpec reports its name
package eval

func (s *derivedSpec) SpecName() string { return s.name }
