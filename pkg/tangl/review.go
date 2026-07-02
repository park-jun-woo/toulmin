//ff:func feature=tangl type=engine control=sequence
//ff:what Review — escalates a failed-compensation pass to a human-review error and logs it
package tangl

import (
	"fmt"
	"os"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

// Review builds a *ReviewError from the original failure (cause) and the
// compensation failure (comp), logs it to stderr (spec: never fail silently),
// and returns it. ctx is accepted for parity with the rest of the runtime API
// and future audit-logging hooks.
func Review(ctx toulmin.Context, cause, comp error) error {
	err := &ReviewError{Cause: cause, Compensate: comp}
	fmt.Fprintln(os.Stderr, "tangl:", err.Error())
	return err
}
