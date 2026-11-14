//ff:type feature=engine type=model
//ff:what Spec — rule judgment criteria interface
package toulmin

// Spec defines judgment criteria for a rule.
// Each Spec is a normalized semantic unit that can be composed with others via With().
type Spec interface {
	SpecName() string
	Validate() error
}
