//ff:func feature=engine type=model control=sequence
//ff:what Validate — rejects a list Spec with an empty name
package toulmin

import "fmt"

// Validate rejects a list Spec with an empty name.
func (s StringListSpec) Validate() error {
	if s.name == "" {
		return fmt.Errorf("string list Spec name is empty")
	}
	return nil
}
