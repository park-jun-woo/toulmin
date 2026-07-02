//ff:func feature=feature type=rule control=sequence
//ff:what TestIsInternalStaff_NonMapAttributes — covers non-map attributes branch of IsInternalStaff
package feature

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestIsInternalStaff_NonMapAttributes(t *testing.T) {
	ctx := toulmin.NewContext()
	ctx.Set("attributes", "not-a-map")
	got, _ := IsInternalStaff(ctx, nil)
	if got != false {
		t.Errorf("got %v, want false", got)
	}
}
