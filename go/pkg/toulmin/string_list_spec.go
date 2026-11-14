//ff:type feature=engine type=model
//ff:what StringListSpec — named immutable string-list judgment criterion
package toulmin

// StringListSpec is the built-in Spec produced by syntax-directed list definitions.
type StringListSpec struct {
	name  string
	items []string
}
