//ff:type feature=tangl type=model
//ff:what Compare — a field/operator/value comparison expression
package ast

// Compare is a leaf condition: `<field> <op> <value>`.
// Op is one of: "is empty", "is not empty", "equals", "is in",
// "is greater than", "is less than", "is at most", "is at least".
type Compare struct {
	Field string `json:"field"`
	Op    string `json:"op"`
	Value string `json:"value,omitempty"`
}
