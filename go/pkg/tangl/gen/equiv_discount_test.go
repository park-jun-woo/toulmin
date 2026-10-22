//ff:type feature=tangl type=codegen
//ff:what equivDiscount — cross-test Spec matched by the `spec discount` reference
package gen

// equivDiscount is the cross-test Spec the “ spec `discount`.`minPrice` “
// reference resolves (SpecName "discount"); the generated runner defines an
// identical copy.
type equivDiscount struct {
	MinPrice float64
}
