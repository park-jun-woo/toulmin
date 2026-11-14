//ff:type feature=engine type=model
//ff:what ruleIDDiscountSpec — nested spec used as pointer field to verify JSON value (not address) drives the rule identifier
package toulmin

// ruleIDDiscountSpec is a nested spec used as a pointer field to verify that
// JSON value serialization (not address) drives the rule identifier.
type ruleIDDiscountSpec struct {
	Name string
	Rate float64
}
