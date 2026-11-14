//ff:func feature=policy type=model control=sequence
//ff:what Validate — checks that HeaderSpec header is non-empty
package policy

import "fmt"

// Validate checks that header is non-empty.
func (b *HeaderSpec) Validate() error {
	if b.Header == "" {
		return fmt.Errorf("HeaderSpec: header is required")
	}
	return nil
}
