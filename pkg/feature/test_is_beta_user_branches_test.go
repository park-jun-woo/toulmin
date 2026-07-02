//ff:func feature=feature type=rule control=sequence
//ff:what TestIsBetaUser_NonMapAttributes — covers non-map attributes branch of IsBetaUser
package feature

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestIsBetaUser_NonMapAttributes(t *testing.T) {
	ctx := toulmin.NewContext()
	ctx.Set("attributes", "not-a-map")
	got, _ := IsBetaUser(ctx, nil)
	if got != false {
		t.Errorf("got %v, want false", got)
	}
}
