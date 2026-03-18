//ff:func feature=engine type=engine control=sequence
//ff:what NewEngine — creates a new Engine instance
package toulmin

// NewEngine returns an empty Engine ready for rule registration.
func NewEngine() *Engine {
	return &Engine{}
}
