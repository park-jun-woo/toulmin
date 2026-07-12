//ff:type feature=tangl type=codegen
//ff:what valCode — emitted statements and expressions resolving one operand value
package gen

// valCode is the emitted code that resolves one operand inside a generated
// leaf closure: Setup statements (possibly empty) followed by the V (value),
// OK (resolution flag), and Actual (witness Fact Actual) expressions. OK and
// Actual are empty when the caller did not request them.
type valCode struct {
	Setup  string
	V      string
	OK     string
	Actual string
}
