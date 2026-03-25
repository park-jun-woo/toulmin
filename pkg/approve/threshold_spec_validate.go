//ff:func feature=approve type=model control=sequence
//ff:what Validate — checks that ThresholdSpec max is positive
package approve

import "fmt"

// Validate checks that max is positive.
func (b *ThresholdSpec) Validate() error {
	if b.Max <= 0 {
		return fmt.Errorf("ThresholdSpec: max must be positive")
	}
	return nil
}
