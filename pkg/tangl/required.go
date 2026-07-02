//ff:func feature=tangl type=engine control=iteration dimension=1
//ff:what Required — verifies that every named context field is present and non-nil
package tangl

import (
	"errors"
	"fmt"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

// ErrRequired is returned (wrapped) when a required context field is missing or nil.
var ErrRequired = errors.New("tangl: required field missing")

// Required checks that every field in fields is present in ctx and non-nil.
// It returns the first missing field wrapped in ErrRequired, or nil if all are present.
func Required(ctx toulmin.Context, fields ...string) error {
	for _, f := range fields {
		v, ok := ctx.Get(f)
		if !ok || v == nil {
			return fmt.Errorf("%w: %s", ErrRequired, f)
		}
	}
	return nil
}
