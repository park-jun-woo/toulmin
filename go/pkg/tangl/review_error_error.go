//ff:func feature=tangl type=model control=sequence
//ff:what ReviewError.Error — human-readable message naming both the cause and the compensation failure
package tangl

import "fmt"

// Error implements the error interface.
func (r *ReviewError) Error() string {
	return fmt.Sprintf("tangl: REVIEW required — cause: %v; compensation failed: %v", r.Cause, r.Compensate)
}
