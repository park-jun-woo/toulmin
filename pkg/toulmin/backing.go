//ff:type feature=engine type=interface
//ff:what Backing — interface for rule judgment criteria (data only, no func fields)
package toulmin

// Backing defines the interface for rule judgment criteria.
// Implementations must be pure data structs (no func fields).
type Backing interface {
	BackingName() string // type identifier for YAML, error messages, logging
	Validate() error     // domain validation — required fields, value ranges
}
