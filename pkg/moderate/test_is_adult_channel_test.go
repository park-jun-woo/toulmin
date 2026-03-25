//ff:func feature=moderate type=rule control=sequence
//ff:what TestIsAdultChannel — tests IsAdultChannel rule
package moderate

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestIsAdultChannel(t *testing.T) {
	ctx := toulmin.NewContext()
	ctx.Set("channel", &Channel{AgeGated: true})
	got, _ := IsAdultChannel(ctx, nil)
	if !got {
		t.Error("expected true for age-gated channel")
	}
}
