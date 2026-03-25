//ff:func feature=approve type=model control=sequence
//ff:what Validate — checks that ThresholdBacking max is positive
package approve

import "fmt"

// Validate checks that max is positive.
func (b *ThresholdBacking) Validate() error {
	if b.Max <= 0 {
		return fmt.Errorf("ThresholdBacking: max must be positive")
	}
	return nil
}
