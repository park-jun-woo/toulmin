//ff:func feature=policy type=model control=sequence
//ff:what Validate — checks that HeaderBacking header is non-empty
package policy

import "fmt"

// Validate checks that header is non-empty.
func (b *HeaderBacking) Validate() error {
	if b.Header == "" {
		return fmt.Errorf("HeaderBacking: header is required")
	}
	return nil
}
