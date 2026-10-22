//ff:type feature=engine type=model
//ff:what ruleIDMemberSpec — spec carrying a nested pointer field to exercise the address stability regression (case d)
package toulmin

// ruleIDMemberSpec carries a nested pointer field to exercise the address
// stability regression (case d).
type ruleIDMemberSpec struct {
	Level    string
	Discount *ruleIDDiscountSpec
}
